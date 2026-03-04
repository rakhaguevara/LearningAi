package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AIImageService calls DashScope `qwen-image-plus` API to create illustrations
// and uploads the results to Alibaba Cloud OSS.
type AIImageService struct {
	apiKey      string
	imageModel  string
	log         *zap.Logger
	client      *http.Client
	ossEndpoint string
	ossBucket   string
	ossKeyID    string
	ossSecret   string
}

// NewAIImageService creates a production-grade image generator and OSS uploader.
func NewAIImageService(apiKey, imageModel, ossEndpoint, ossBucket, ossKeyID, ossSecret string, log *zap.Logger) *AIImageService {
	return &AIImageService{
		apiKey:      apiKey,
		imageModel:  imageModel,
		log:         log,
		client:      &http.Client{Timeout: 60 * time.Second},
		ossEndpoint: ossEndpoint,
		ossBucket:   ossBucket,
		ossKeyID:    ossKeyID,
		ossSecret:   ossSecret,
	}
}

// GenerateImage submits a text-to-image task to qwen-image-plus, downloads the image, and uploads it to OSS.
// Returns the public OSS URL of the image.
func (s *AIImageService) GenerateImage(ctx context.Context, prompt string) (*ImageResult, error) {
	input := ImageGenerationInput{
		FinalPrompt:    prompt,
		NegativePrompt: DefaultNegativePrompt,
		Style:          "flat_educational",
	}
	return s.GenerateImageFromInput(ctx, input)
}

// GenerateImageFromInput is the preferred entry-point. It uses a fully-prepared
// ImageGenerationInput (enhanced prompt + negative prompt + style).
func (s *AIImageService) GenerateImageFromInput(ctx context.Context, input ImageGenerationInput) (*ImageResult, error) {
	if s.apiKey == "" {
		s.log.Warn("image_generation_skipped_no_api_key")
		return &ImageResult{Fallback: true, FallbackMsg: "API key not configured in AIImageService"}, nil
	}

	start := time.Now()
	s.log.Info("image_generation_started",
		zap.String("image_model", s.imageModel),
		zap.String("style", input.Style),
		zap.String("image_prompt_preview", truncate(input.FinalPrompt, 150)),
		zap.String("negative_prompt_preview", truncate(input.NegativePrompt, 80)),
	)

	imageURL, err := s.callDashscopeImageAPI(ctx, input.FinalPrompt, input.NegativePrompt)
	if err != nil {
		durationMs := int(time.Since(start).Milliseconds())
		s.log.Error("image_generation_failed",
			zap.String("image_error", err.Error()),
			zap.Int("generation_duration", durationMs),
		)
		// Return fallback result but don't return error - allows orchestrator to continue
		return &ImageResult{Fallback: true, FallbackMsg: err.Error(), DurationMs: durationMs}, nil
	}

	s.log.Info("dashscope_image_generated_successfully",
		zap.String("generated_image_url", imageURL),
		zap.Int("url_length", len(imageURL)),
	)

	// Try to upload to OSS if configured, but ALWAYS return the DashScope URL even if OSS fails
	if s.ossEndpoint != "" && s.ossBucket != "" && s.ossKeyID != "" && s.ossSecret != "" {
		ossURL, ossErr := s.downloadAndUploadToOSS(ctx, imageURL)
		if ossErr != nil {
			s.log.Warn("oss_upload_failed_returning_dashscope_url",
				zap.String("dashscope_url", imageURL),
				zap.String("oss_error", ossErr.Error()),
			)
			// Return the original DashScope URL which is valid for 24 hours
			return &ImageResult{
				ImageURL:    imageURL,
				Model:       s.imageModel,
				GeneratedAt: time.Now(),
				DurationMs:  int(time.Since(start).Milliseconds()),
				Fallback:    false,
			}, nil
		}

		durationMs := int(time.Since(start).Milliseconds())
		s.log.Info("image_uploaded_to_oss_successfully",
			zap.String("oss_url", ossURL),
			zap.String("image_model", s.imageModel),
			zap.Int("generation_duration", durationMs),
		)

		return &ImageResult{
			ImageURL:    ossURL,
			Model:       s.imageModel,
			GeneratedAt: time.Now(),
			DurationMs:  durationMs,
			Fallback:    false,
		}, nil
	}

	// OSS not configured - return DashScope URL directly
	durationMs := int(time.Since(start).Milliseconds())
	s.log.Info("image_generation_complete_using_dashscope_url",
		zap.String("dashscope_url", imageURL),
		zap.String("image_model", s.imageModel),
		zap.Int("generation_duration", durationMs),
	)

	return &ImageResult{
		ImageURL:    imageURL,
		Model:       s.imageModel,
		GeneratedAt: time.Now(),
		DurationMs:  durationMs,
		Fallback:    false,
	}, nil
}

func (s *AIImageService) callDashscopeImageAPI(ctx context.Context, prompt string, negativePrompt string) (string, error) {
	url := "https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis"

	type taskInput struct {
		Prompt         string `json:"prompt"`
		NegativePrompt string `json:"negative_prompt,omitempty"`
	}
	type taskParams struct {
		Size         string `json:"size"`
		N            int    `json:"n"`
		PromptExtend bool   `json:"prompt_extend"`
		Watermark    bool   `json:"watermark"`
	}
	type taskBody struct {
		Model      string     `json:"model"`
		Input      taskInput  `json:"input"`
		Parameters taskParams `json:"parameters"`
	}

	body := taskBody{
		Model: s.imageModel,
		Input: taskInput{
			Prompt:         prompt,
			NegativePrompt: negativePrompt,
		},
		Parameters: taskParams{
			Size:         "1664*928",
			N:            1,
			PromptExtend: true,
			Watermark:    false,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshalling image request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("building image request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	s.log.Info("submitting_image_request_to_dashscope",
		zap.String("url", url),
		zap.String("model", s.imageModel),
		zap.String("prompt_preview", truncate(prompt, 100)),
	)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("image generation http call: %w", err)
	}
	defer resp.Body.Close()

	bodyBytesRes, _ := io.ReadAll(resp.Body)
	s.log.Debug("dashscope_image_raw_response", zap.String("raw_body", string(bodyBytesRes)))

	if resp.StatusCode != http.StatusOK {
		s.log.Error("dashscope_api_http_error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", truncate(string(bodyBytesRes), 300)),
		)
		return "", fmt.Errorf("image generation HTTP %d: %s", resp.StatusCode, truncate(string(bodyBytesRes), 200))
	}

	var result struct {
		Output struct {
			TaskID     string `json:"task_id"`
			TaskStatus string `json:"task_status"`
			Results    []struct {
				URL string `json:"url"`
			} `json:"results"`
		} `json:"output"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(bodyBytesRes, &result); err != nil {
		return "", fmt.Errorf("decoding image response: %w (raw: %s)", err, truncate(string(bodyBytesRes), 200))
	}

	// Check for API-level errors
	if result.Code != "" || result.Message != "" {
		s.log.Error("dashscope_api_error_returned",
			zap.String("code", result.Code),
			zap.String("message", result.Message),
			zap.String("raw_response", string(bodyBytesRes)),
		)
		return "", fmt.Errorf("API error: %s - %s", result.Code, result.Message)
	}

	// Validate task ID
	if result.Output.TaskID == "" {
		s.log.Error("no_task_id_in_response",
			zap.String("raw_response", string(bodyBytesRes)),
		)
		return "", fmt.Errorf("no task_id in response")
	}

	s.log.Info("image_task_submitted_successfully",
		zap.String("task_id", result.Output.TaskID),
		zap.String("task_status", result.Output.TaskStatus),
	)

	// Poll for task completion
	imageURL, err := s.pollTask(ctx, result.Output.TaskID)
	if err != nil {
		s.log.Error("polling_task_failed",
			zap.String("task_id", result.Output.TaskID),
			zap.Error(err),
		)
		return "", fmt.Errorf("polling task failed: %w", err)
	}

	if imageURL == "" {
		s.log.Error("empty_image_url_in_results",
			zap.String("task_id", result.Output.TaskID),
		)
		return "", fmt.Errorf("empty image URL in results")
	}

	s.log.Info("image_url_extracted_from_polling",
		zap.String("image_url", imageURL),
		zap.Int("url_length", len(imageURL)),
	)

	return imageURL, nil
}

// pollTask polls the DashScope task endpoint until completion
func (s *AIImageService) pollTask(ctx context.Context, taskID string) (string, error) {
	deadline := time.Now().Add(65 * time.Second)
	queryURL := "https://dashscope-intl.aliyuncs.com/api/v1/tasks/" + taskID

	for {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("task timed out after 65s (task_id: %s)", taskID)
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("context cancelled while polling task %s", taskID)
		case <-time.After(3 * time.Second):
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
		if err != nil {
			s.log.Warn("task poll build request error", zap.Error(err))
			continue
		}
		req.Header.Set("Authorization", "Bearer "+s.apiKey)

		resp, err := s.client.Do(req)
		if err != nil {
			s.log.Warn("task poll http error", zap.String("task_id", taskID), zap.Error(err))
			continue
		}

		data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		resp.Body.Close()

		s.log.Debug("task poll response",
			zap.String("task_id", taskID),
			zap.Int("status", resp.StatusCode),
			zap.String("body", truncate(string(data), 400)),
		)

		var result struct {
			Output struct {
				TaskStatus string `json:"task_status"`
				Results    []struct {
					URL  string `json:"url"`
					Code string `json:"code"`
				} `json:"results"`
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"output"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}

		if err := json.Unmarshal(data, &result); err != nil {
			s.log.Warn("task poll unmarshal error", zap.Error(err))
			continue
		}

		s.log.Debug("task poll status",
			zap.String("task_id", taskID),
			zap.String("task_status", result.Output.TaskStatus),
		)

		switch result.Output.TaskStatus {
		case "SUCCEEDED":
			if len(result.Output.Results) == 0 {
				return "", fmt.Errorf("task SUCCEEDED but no results returned (task_id: %s)", taskID)
			}
			imageURL := strings.TrimSpace(result.Output.Results[0].URL)
			if imageURL == "" {
				return "", fmt.Errorf("task SUCCEEDED but image URL is empty (task_id: %s)", taskID)
			}
			s.log.Info("task_succeeded_with_url",
				zap.String("task_id", taskID),
				zap.String("image_url", imageURL),
			)
			return imageURL, nil

		case "FAILED":
			return "", fmt.Errorf("task FAILED: %s (task_id: %s)", result.Output.Message, taskID)

		default:
			// PENDING / RUNNING — keep polling
			s.log.Debug("task still running", zap.String("status", result.Output.TaskStatus))
		}
	}
}

func (s *AIImageService) downloadAndUploadToOSS(ctx context.Context, imageURL string) (string, error) {
	if s.ossEndpoint == "" || s.ossBucket == "" || s.ossKeyID == "" || s.ossSecret == "" {
		return "", fmt.Errorf("OSS credentials not configured")
	}

	// 1. Download Image
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("create download request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download image HTTP %d", resp.StatusCode)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read image: %w", err)
	}

	// 2. Upload to OSS
	client, err := oss.New(s.ossEndpoint, s.ossKeyID, s.ossSecret)
	if err != nil {
		return "", fmt.Errorf("OSS client init: %w", err)
	}

	bucket, err := client.Bucket(s.ossBucket)
	if err != nil {
		return "", fmt.Errorf("OSS bucket init: %w", err)
	}

	objectKey := fmt.Sprintf("illustrations/img_%d_%s.png", time.Now().UnixNano(), uuid.New().String()[:8])

	err = bucket.PutObject(objectKey, bytes.NewReader(imageBytes))
	if err != nil {
		return "", fmt.Errorf("OSS put object: %w", err)
	}

	// 3. Construct Public URL
	// E.g., https://ailearn-assets.oss-cn-hangzhou.aliyuncs.com/...
	publicURL := fmt.Sprintf("https://%s.%s/%s", s.ossBucket, s.ossEndpoint, objectKey)
	return publicURL, nil
}
