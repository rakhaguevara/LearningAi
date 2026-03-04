package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	if s.apiKey == "" {
		return nil, fmt.Errorf("API key not configured in AIImageService")
	}

	start := time.Now()
	s.log.Info("image_generation_triggered",
		zap.String("image_model", s.imageModel),
		zap.String("image_prompt", truncate(prompt, 120)),
	)

	imageURL, err := s.callDashscopeImageAPI(ctx, prompt)
	if err != nil {
		durationMs := int(time.Since(start).Milliseconds())
		s.log.Warn("image_generation_failed",
			zap.String("image_error", err.Error()),
			zap.Int("generation_duration", durationMs),
		)
		return &ImageResult{Fallback: true, FallbackMsg: err.Error()}, err
	}

	ossURL, err := s.downloadAndUploadToOSS(ctx, imageURL)
	if err != nil {
		durationMs := int(time.Since(start).Milliseconds())
		s.log.Warn("image_persistence_failed",
			zap.String("image_error", err.Error()),
			zap.Int("generation_duration", durationMs),
		)
		// Still return the generated URL even if upload failed, it's valid for 24 hours.
		return &ImageResult{
			ImageURL:    imageURL,
			Model:       s.imageModel,
			GeneratedAt: time.Now(),
			DurationMs:  int(time.Since(start).Milliseconds()),
			Fallback:    false,
		}, nil
	}

	durationMs := int(time.Since(start).Milliseconds())
	s.log.Info("image_generation_complete",
		zap.String("image_url", ossURL),
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

func (s *AIImageService) callDashscopeImageAPI(ctx context.Context, prompt string) (string, error) {
	url := "https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis"

	type taskInput struct {
		Prompt string `json:"prompt"`
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
		Input: taskInput{Prompt: prompt},
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

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("image generation http call: %w", err)
	}
	defer resp.Body.Close()

	bodyBytesRes, _ := io.ReadAll(resp.Body)
	s.log.Debug("dashscope_image_raw_response", zap.String("raw_body", string(bodyBytesRes)))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image generation http %d: %s", resp.StatusCode, string(bodyBytesRes))
	}

	var result struct {
		Output struct {
			TaskStatus string `json:"task_status"`
			Results    []struct {
				URL string `json:"url"`
			} `json:"results"`
		} `json:"output"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(bodyBytesRes, &result); err != nil {
		return "", fmt.Errorf("decoding image response: %w", err)
	}

	if result.Code != "" || result.Message != "" {
		s.log.Error("dashscope_api_error_returned",
			zap.String("code", result.Code),
			zap.String("message", result.Message),
			zap.String("raw_response", string(bodyBytesRes)),
		)
		return "", fmt.Errorf("API error: %s - %s", result.Code, result.Message)
	}

	if len(result.Output.Results) == 0 {
		return "", fmt.Errorf("no visual result in API response")
	}

	return result.Output.Results[0].URL, nil
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
