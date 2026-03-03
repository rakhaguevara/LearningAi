package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// TTSRequest is the input to the TTS service.
type TTSRequest struct {
	Text     string  `json:"text"`
	Voice    string  `json:"voice,omitempty"`     // e.g. "longxiaochun"
	Format   string  `json:"format,omitempty"`    // mp3, wav, pcm
	SampleHz int     `json:"sample_hz,omitempty"` // 16000, 24000
	Speed    float64 `json:"speed,omitempty"`     // 0.5–2.0
}

// TTSResult is returned after audio generation.
type TTSResult struct {
	FilePath  string    `json:"file_path"`
	FileName  string    `json:"file_name"`
	Format    string    `json:"format"`
	Duration  float64   `json:"duration_sec"`
	ByteSize  int       `json:"byte_size"`
	CreatedAt time.Time `json:"created_at"`
}

// TTSService calls Alibaba CosyVoice / DashScope TTS API.
type TTSService struct {
	apiKey    string
	baseURL   string
	outputDir string
	log       *zap.Logger
	client    *http.Client
}

func NewTTSService(apiKey, baseURL, outputDir string, log *zap.Logger) *TTSService {
	if outputDir == "" {
		outputDir = "/tmp/ailearn/audio"
	}
	return &TTSService{
		apiKey:    apiKey,
		baseURL:   baseURL,
		outputDir: outputDir,
		log:       log,
		client:    &http.Client{Timeout: 120 * time.Second},
	}
}

// GenerateAudio converts text to speech and writes the audio file to disk.
func (t *TTSService) GenerateAudio(ctx context.Context, userID string, req TTSRequest) (*TTSResult, error) {
	if t.apiKey == "" {
		return nil, fmt.Errorf("QWEN_API_KEY is not configured for TTS")
	}
	if len(req.Text) > 5000 {
		req.Text = req.Text[:5000] // DashScope CosyVoice limit
	}

	// Apply defaults
	voice := coalesce(req.Voice, "longxiaochun")
	format := coalesce(req.Format, "mp3")
	sampleRate := req.SampleHz
	if sampleRate == 0 {
		sampleRate = 24000
	}
	speed := req.Speed
	if speed == 0 {
		speed = 1.0
	}

	// Build request for DashScope CosyVoice TTS
	type ttsInput struct {
		Text string `json:"text"`
	}
	type ttsParams struct {
		TextType   string  `json:"text_type"`
		Voice      string  `json:"voice"`
		Format     string  `json:"format"`
		SampleRate int     `json:"sample_rate"`
		Volume     int     `json:"volume"`
		Rate       float64 `json:"rate"`
	}
	type ttsBody struct {
		Model      string    `json:"model"`
		Input      ttsInput  `json:"input"`
		Parameters ttsParams `json:"parameters"`
	}

	body := ttsBody{
		Model: "cosyvoice-v1",
		Input: ttsInput{Text: req.Text},
		Parameters: ttsParams{
			TextType:   "PlainText",
			Voice:      voice,
			Format:     format,
			SampleRate: sampleRate,
			Volume:     50,
			Rate:       speed,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling TTS request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		t.baseURL+"/services/aigc/text2audio/synthesis",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("creating TTS request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+t.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-DashScope-DataInspection", "disable")

	start := time.Now()
	httpResp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing TTS request: %w", err)
	}
	defer httpResp.Body.Close()

	audioData, err := io.ReadAll(io.LimitReader(httpResp.Body, 32*1024*1024)) // 32 MB cap
	if err != nil {
		return nil, fmt.Errorf("reading TTS audio: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TTS http %d: %s", httpResp.StatusCode, string(audioData))
	}

	// Persist to disk
	userDir := filepath.Join(t.outputDir, sanitisePathSegment(userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("creating audio directory: %w", err)
	}

	fileName := fmt.Sprintf("audio_%d.%s", time.Now().UnixMilli(), format)
	filePath := filepath.Join(userDir, fileName)

	if err := os.WriteFile(filePath, audioData, 0644); err != nil {
		return nil, fmt.Errorf("writing audio file: %w", err)
	}

	latency := time.Since(start)
	t.log.Info("TTS generated",
		zap.String("user_id", userID),
		zap.String("file", filePath),
		zap.Int("bytes", len(audioData)),
		zap.Duration("latency", latency),
	)

	// Rough duration estimate: ~150 words/min, 5 chars/word avg
	wordCount := len(req.Text) / 5
	durationSec := float64(wordCount) / 150.0 * 60.0

	return &TTSResult{
		FilePath:  filePath,
		FileName:  fileName,
		Format:    format,
		Duration:  durationSec,
		ByteSize:  len(audioData),
		CreatedAt: time.Now(),
	}, nil
}
