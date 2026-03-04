package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ──────────────────────────────────────────────────────────────────────────────
// OpenAI-compatible API wire types (DashScope intl compatible-mode)
// ──────────────────────────────────────────────────────────────────────────────

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIChatRequest is the OpenAI-compatible request format.
type openAIChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
}

// openAIChatResponse is the OpenAI-compatible response format.
type openAIChatResponse struct {
	Choices []struct {
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	ID    string `json:"id"`
	Model string `json:"model"`
	// Error fields (returned with HTTP 4xx/5xx)
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Domain types
// ──────────────────────────────────────────────────────────────────────────────

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
	// Optional overrides
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Content      string `json:"content"`
	TokensUsed   int    `json:"tokens_used"`
	LatencyMs    int    `json:"latency_ms"`
	FinishReason string `json:"finish_reason"`
}

// ──────────────────────────────────────────────────────────────────────────────
// QwenClient — production-grade HTTP client with retry / backoff
// ──────────────────────────────────────────────────────────────────────────────

type QwenClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	log        *zap.Logger
	maxRetries int
}

func NewQwenClient(apiKey, baseURL, model string, log *zap.Logger) *QwenClient {
	return &QwenClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
		log:        log,
		maxRetries: 3,
	}
}

// GenerateChatCompletion sends a multi-turn chat to DashScope (OpenAI-compatible)
// and returns the assistant reply, token usage, and latency.
func (c *QwenClient) GenerateChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if c.apiKey == "" || c.apiKey == "your-qwen-api-key" {
		return nil, fmt.Errorf("QWEN_API_KEY is not configured — set a real DashScope API key in your .env file")
	}

	maxTokens := 2048
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}
	temperature := 0.7
	if req.Temperature > 0 {
		temperature = req.Temperature
	}

	body := openAIChatRequest{
		Model:       c.model,
		Messages:    req.Messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        0.8,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	start := time.Now()

	var (
		resp    *openAIChatResponse
		attempt int
	)

	for attempt = 1; attempt <= c.maxRetries; attempt++ {
		resp, err = c.doRequest(ctx, bodyBytes)
		if err == nil {
			break
		}

		// Don't retry on authentication / bad-request errors — they won't recover.
		errStr := err.Error()
		if strings.Contains(errStr, "qwen http 401") ||
			strings.Contains(errStr, "qwen http 403") ||
			strings.Contains(errStr, "qwen http 400") {
			return nil, err // fail fast
		}

		if attempt == c.maxRetries {
			return nil, fmt.Errorf("qwen api failed after %d attempts: %w", c.maxRetries, err)
		}

		backoff := time.Duration(math.Pow(2, float64(attempt))) * 500 * time.Millisecond
		c.log.Warn("qwen request failed, retrying",
			zap.Int("attempt", attempt),
			zap.Duration("backoff", backoff),
			zap.Error(err),
		)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
	}

	latency := int(time.Since(start).Milliseconds())

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("qwen returned no choices")
	}

	choice := resp.Choices[0]

	c.log.Info("qwen chat completed",
		zap.Int("tokens_total", resp.Usage.TotalTokens),
		zap.Int("latency_ms", latency),
		zap.Int("attempt", attempt),
		zap.String("finish_reason", choice.FinishReason),
	)

	return &ChatResponse{
		Content:      choice.Message.Content,
		TokensUsed:   resp.Usage.TotalTokens,
		LatencyMs:    latency,
		FinishReason: choice.FinishReason,
	}, nil
}

func (c *QwenClient) doRequest(ctx context.Context, body []byte) (*openAIChatResponse, error) {
	// OpenAI-compatible endpoint on DashScope international
	url := c.baseURL + "/chat/completions"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing http request: %w", err)
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(httpResp.Body, 4*1024*1024)) // 4 MB cap
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("qwen http %d: %s", httpResp.StatusCode, string(respBytes))
	}

	var result openAIChatResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("decoding qwen response: %w", err)
	}

	// OpenAI-compatible error field
	if result.Error != nil {
		return nil, fmt.Errorf("qwen api error %s: %s", result.Error.Code, result.Error.Message)
	}

	return &result, nil
}
