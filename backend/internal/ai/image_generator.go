package ai

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// ──────────────────────────────────────────────────────────────────────────────
// ImageGenerator — DashScope Wanx Text-to-Image
// ──────────────────────────────────────────────────────────────────────────────

const (
	defaultImageModel   = "wanx-v1"
	maxImagePollSeconds = 65
)

// ImagePromptStyle enumerates which image style prefix to use.
type ImagePromptStyle string

const (
	ImageStyleAnime  ImagePromptStyle = "anime"
	ImageStyleSports ImagePromptStyle = "sports"
)

// ImageResult is returned by GenerateImage.
type ImageResult struct {
	ImageURL    string    `json:"image_url"`
	Model       string    `json:"model"`
	GeneratedAt time.Time `json:"generated_at"`
	DurationMs  int       `json:"duration_ms"`
	Fallback    bool      `json:"fallback"` // true if no image could be generated
	FallbackMsg string    `json:"fallback_msg,omitempty"`
	LocalPath   string    `json:"local_path,omitempty"` // local file path if downloaded
}

// ImageGenerator calls DashScope Wanx API to create illustrations.
type ImageGenerator struct {
	apiKey     string
	baseURL    string // e.g. https://dashscope-intl.aliyuncs.com/api/v1
	model      string
	log        *zap.Logger
	client     *http.Client
	storageDir string // local directory to store generated images
}

// NewImageGenerator creates a production-grade image generator.
func NewImageGenerator(apiKey, baseURL, storageDir string, log *zap.Logger) *ImageGenerator {
	model := os.Getenv("QWEN_IMAGE_MODEL")
	if model == "" {
		model = defaultImageModel
	}

	// Force the standard DashScope endpoint for image generation
	// because `wanx-v1` "Model not exist" occurs on the `-intl` subdomain.
	imageURL := "https://dashscope.aliyuncs.com/api/v1"

	// Ensure storage directory exists
	if storageDir == "" {
		storageDir = "/tmp/ailearn/generated-images"
	}
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Warn("failed to create image storage directory", zap.Error(err))
		storageDir = "" // disable local storage if we can't create the dir
	}

	return &ImageGenerator{
		apiKey:     apiKey,
		baseURL:    imageURL,
		model:      model,
		log:        log,
		client:     &http.Client{Timeout: 95 * time.Second},
		storageDir: storageDir,
	}
}

// BuildImagePrompt prepends a style-specific wrapper around the raw scene prompt.
func BuildImagePrompt(style ImagePromptStyle, scenePrompt string) string {
	switch style {
	case ImageStyleAnime:
		return fmt.Sprintf(
			"Anime style illustration, cinematic lighting, dynamic action scene, vibrant colors, detailed background, educational theme, %s",
			scenePrompt,
		)
	case ImageStyleSports:
		return fmt.Sprintf(
			"Dynamic sports action scene, motion blur, energetic lighting, educational visualization, %s",
			scenePrompt,
		)
	default:
		return scenePrompt
	}
}

// GenerateImage submits a text-to-image task and polls for completion.
// Always returns an ImageResult — never panics or returns nil.
// On failure sets Fallback=true and logs the error.
func (g *ImageGenerator) GenerateImage(ctx context.Context, prompt string, style ImagePromptStyle) *ImageResult {
	if g.apiKey == "" {
		g.log.Warn("image generation skipped — no API key")
		return &ImageResult{Fallback: true, FallbackMsg: "API key not configured"}
	}

	start := time.Now()
	g.log.Info("image_generation_triggered",
		zap.Bool("image_triggered", true),
		zap.String("model", g.model),
		zap.String("image_prompt", truncate(prompt, 120)),
		zap.String("style", string(style)),
	)

	taskID, err := g.submitTask(ctx, prompt, style)
	if err != nil {
		durationMs := int(time.Since(start).Milliseconds())
		g.log.Warn("image_generation_submit_failed",
			zap.Bool("image_triggered", true),
			zap.String("image_error", err.Error()),
			zap.Int("image_generation_duration", durationMs),
		)
		// Retry once
		g.log.Info("image_generation_retry_attempt")
		taskID, err = g.submitTask(ctx, prompt, style)
		if err != nil {
			durationMs = int(time.Since(start).Milliseconds())
			g.log.Warn("image_generation_retry_failed",
				zap.String("image_error", err.Error()),
				zap.Int("image_generation_duration", durationMs),
			)
			return &ImageResult{Fallback: true, FallbackMsg: fmt.Sprintf("submit failed after retry: %v", err)}
		}
	}

	g.log.Info("image task submitted", zap.String("task_id", taskID))

	imageURL, err := g.pollTask(ctx, taskID)
	durationMs := int(time.Since(start).Milliseconds())

	if err != nil {
		g.log.Warn("image_generation_poll_failed",
			zap.String("task_id", taskID),
			zap.String("image_error", err.Error()),
			zap.Int("image_generation_duration", durationMs),
		)
		return &ImageResult{Fallback: true, FallbackMsg: fmt.Sprintf("poll failed: %v", err)}
	}

	g.log.Info("image_generation_complete",
		zap.Bool("image_triggered", true),
		zap.String("task_id", taskID),
		zap.String("image_url", imageURL),
		zap.Int("image_generation_duration", durationMs),
		zap.Bool("fallback_triggered", false),
	)

	// Download and save image locally if storage is enabled
	localURL := ""
	if g.storageDir != "" && imageURL != "" {
		if localPath, err := g.downloadAndSaveImage(ctx, imageURL); err != nil {
			g.log.Warn("image_download_failed", zap.Error(err))
		} else {
			localURL = localPath
			g.log.Info("image_downloaded_successfully", zap.String("local_path", localPath))
		}
	}

	return &ImageResult{
		ImageURL:    localURL, // Return local path if successful, otherwise remote URL
		Model:       g.model,
		GeneratedAt: time.Now(),
		DurationMs:  durationMs,
		Fallback:    false,
		LocalPath:   localURL,
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Internal: submit task
// ──────────────────────────────────────────────────────────────────────────────

type wanxTaskBody struct {
	Model      string         `json:"model"`
	Input      wanxTaskInput  `json:"input"`
	Parameters wanxTaskParams `json:"parameters"`
}

type wanxTaskInput struct {
	Prompt string `json:"prompt"`
}

type wanxTaskParams struct {
	Style string `json:"style"`
	Size  string `json:"size"`
	N     int    `json:"n"`
}

type wanxSubmitResponse struct {
	Output struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
	} `json:"output"`
	Code    string `json:"code"`
	Message string `json:"message"`
	// DashScope often wraps errors differently
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type wanxQueryResponse struct {
	Output struct {
		TaskStatus string `json:"task_status"`
		Results    []struct {
			URL  string `json:"url"`
			Code string `json:"code"`
		} `json:"results"`
		// Error info
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"output"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (g *ImageGenerator) submitTask(ctx context.Context, prompt string, style ImagePromptStyle) (string, error) {
	// Map ImagePromptStyle to DashScope style string
	styleStr := "<anime>"
	if style == ImageStyleSports {
		styleStr = "<anime>" // Wanx only supports anime; sports uses the same style token
	}

	body := wanxTaskBody{
		Model: g.model,
		Input: wanxTaskInput{Prompt: prompt},
		Parameters: wanxTaskParams{
			Style: styleStr,
			Size:  "1024*1024",
			N:     1,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshalling image request: %w", err)
	}

	// Correct DashScope URL for image synthesis
	submitURL := g.baseURL + "/services/aigc/text2image/image-synthesis"
	g.log.Debug("image submit URL", zap.String("url", submitURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, submitURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("building image submit request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("image submit http call: %w", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	g.log.Debug("image submit raw response",
		zap.Int("status", resp.StatusCode),
		zap.String("body", truncate(string(data), 300)),
	)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image submit HTTP %d: %s", resp.StatusCode, string(data))
	}

	var result wanxSubmitResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("decoding submit response: %w (raw: %s)", err, truncate(string(data), 200))
	}

	if result.Output.TaskID == "" {
		msg := result.Message
		if result.Error != nil {
			msg = result.Error.Message
		}
		return "", fmt.Errorf("no task_id in submit response — API message: %s (code: %s)", msg, result.Code)
	}

	return result.Output.TaskID, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Internal: poll task until SUCCEEDED or FAILED
// ──────────────────────────────────────────────────────────────────────────────

func (g *ImageGenerator) pollTask(ctx context.Context, taskID string) (string, error) {
	deadline := time.Now().Add(maxImagePollSeconds * time.Second)

	// Correct DashScope task query URL
	queryURL := g.baseURL + "/tasks/" + taskID

	for {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("image task timed out after %ds (task_id: %s)", maxImagePollSeconds, taskID)
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("context cancelled while polling image task %s", taskID)
		case <-time.After(4 * time.Second):
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
		if err != nil {
			g.log.Warn("image poll build request error", zap.Error(err))
			continue
		}
		req.Header.Set("Authorization", "Bearer "+g.apiKey)

		resp, err := g.client.Do(req)
		if err != nil {
			g.log.Warn("image poll http error", zap.String("task_id", taskID), zap.Error(err))
			continue
		}

		data, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		resp.Body.Close()

		g.log.Debug("image poll response",
			zap.String("task_id", taskID),
			zap.Int("status", resp.StatusCode),
			zap.String("body", truncate(string(data), 400)),
		)

		var result wanxQueryResponse
		if err := json.Unmarshal(data, &result); err != nil {
			g.log.Warn("image poll unmarshal error", zap.Error(err))
			continue
		}

		g.log.Debug("image poll status",
			zap.String("task_id", taskID),
			zap.String("task_status", result.Output.TaskStatus),
		)

		switch result.Output.TaskStatus {
		case "SUCCEEDED":
			if len(result.Output.Results) == 0 {
				return "", fmt.Errorf("task SUCCEEDED but no results returned (task_id: %s)", taskID)
			}
			imageURL := result.Output.Results[0].URL
			if imageURL == "" {
				return "", fmt.Errorf("task SUCCEEDED but image URL is empty (task_id: %s)", taskID)
			}
			return imageURL, nil

		case "FAILED":
			return "", fmt.Errorf("image task FAILED — %s (task_id: %s)", result.Output.Message, taskID)

		default:
			// PENDING / RUNNING — keep polling
			g.log.Debug("image task still running", zap.String("status", result.Output.TaskStatus))
		}
	}
}

// downloadAndSaveImage downloads the image from remote URL and saves it locally.
// Returns the local URL path (e.g., /generated/ai-image-xxxx.png).
func (g *ImageGenerator) downloadAndSaveImage(ctx context.Context, imageURL string) (string, error) {
	if g.storageDir == "" {
		return "", fmt.Errorf("storage directory not configured")
	}

	// Fetch the image
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating download request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("downloading image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Read image data
	imgData, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20)) // Max 20MB
	if err != nil {
		return "", fmt.Errorf("reading image data: %w", err)
	}

	// Determine file extension from Content-Type or default to .png
	contentType := resp.Header.Get("Content-Type")
	ext := ".png"
	switch contentType {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}

	// Generate unique filename
	filename := fmt.Sprintf("ai-image-%s%s", generateUniqueID(), ext)
	filePath := filepath.Join(g.storageDir, filename)

	// Save to disk
	if err := os.WriteFile(filePath, imgData, 0644); err != nil {
		return "", fmt.Errorf("saving image to disk: %w", err)
	}

	g.log.Info("image_saved_locally",
		zap.String("file_path", filePath),
		zap.Int("size_bytes", len(imgData)),
	)

	// Return URL-safe path (relative to server root)
	// This will be served via static file middleware
	return "/generated/" + filename, nil
}

// generateUniqueID creates a random hex string for unique filenames.
func generateUniqueID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
