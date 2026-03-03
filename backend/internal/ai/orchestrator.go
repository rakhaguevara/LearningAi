package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ──────────────────────────────────────────────────────────────────────────────
// Typed schema structs — one per output format
// ──────────────────────────────────────────────────────────────────────────────

type SummaryResponse struct {
	Title                  string   `json:"title"`
	CoreConceptExplanation string   `json:"core_concept_explanation"`
	KeyPoints              []string `json:"key_points"`
	RealWorldExample       string   `json:"real_world_example"`
	ShortConclusion        string   `json:"short_conclusion"`
}

type DetailedResponse struct {
	Title               string   `json:"title"`
	ConceptExplanation  string   `json:"concept_explanation"`
	StepByStepBreakdown []string `json:"step_by_step_breakdown"`
	Example             string   `json:"example"`
	MiniQuiz            []string `json:"mini_quiz"`
}

type AnimeResponse struct {
	EpisodeTitle       string `json:"episode_title"`
	MainCharacter      string `json:"main_character"`
	StoryArc           string `json:"story_arc"`
	PhysicsExplanation string `json:"physics_explanation"`
	VisualScenePrompt  string `json:"visual_scene_prompt"`
}

type SportsResponse struct {
	GameTitle         string `json:"game_title"`
	SportUsed         string `json:"sport_used"`
	PlayBreakdown     string `json:"play_breakdown"`
	CoachingTip       string `json:"coaching_tip"`
	ScoreboardSummary string `json:"scoreboard_summary"`
	VisualScenePrompt string `json:"visual_scene_prompt"`
}

type AcademicResponse struct {
	Title                 string `json:"title"`
	Abstract              string `json:"abstract"`
	TheoreticalBackground string `json:"theoretical_background"`
	Methodology           string `json:"methodology"`
	Conclusion            string `json:"conclusion"`
}

// ──────────────────────────────────────────────────────────────────────────────
// FinalResponse — the fully-parsed, typed response returned to the API handler
// ──────────────────────────────────────────────────────────────────────────────

// FinalResponse is what the orchestrator returns. The API handler converts this
// into the AskResponse JSON that the frontend consumes.
type FinalResponse struct {
	// Core identity
	Mode    OutputFormat `json:"mode"`
	RawJSON string       `json:"raw_json"` // canonical JSON for debugging

	// Typed content — only one of these is non-nil at a time
	Summary  *SummaryResponse  `json:"summary,omitempty"`
	Detailed *DetailedResponse `json:"detailed,omitempty"`
	Anime    *AnimeResponse    `json:"anime,omitempty"`
	Sports   *SportsResponse   `json:"sports,omitempty"`
	Academic *AcademicResponse `json:"academic,omitempty"`

	// Flat convenience fields (set regardless of mode)
	Title           string `json:"title"`
	MainContent     string `json:"main_content"`
	IllustrationURL string `json:"illustration_url"`

	// Metadata
	TokensUsed       int  `json:"tokens_used"`
	OrchestrationMs  int  `json:"orchestration_ms"`
	IsStructuredJSON bool `json:"is_structured_json"`
	JSONParseSuccess bool `json:"json_parse_success"`
}

// ──────────────────────────────────────────────────────────────────────────────
// AIOrchestrator
// ──────────────────────────────────────────────────────────────────────────────

// AIOrchestrator owns the full multimodal pipeline:
//  1. Calls Qwen text model
//  2. Parses JSON into typed structs
//  3. Triggers image generation for visual modes
//  4. Returns a FinalResponse with all data populated
type AIOrchestrator struct {
	textClient *QwenClient
	imageGen   *ImageGenerator
	log        *zap.Logger
}

// NewAIOrchestrator creates the orchestrator.
func NewAIOrchestrator(textClient *QwenClient, imageGen *ImageGenerator, log *zap.Logger) *AIOrchestrator {
	return &AIOrchestrator{
		textClient: textClient,
		imageGen:   imageGen,
		log:        log,
	}
}

// GenerateStructured is the main orchestration entrypoint for structured formats.
// It handles: text gen → JSON parse → retry → image gen → typed response.
func (o *AIOrchestrator) GenerateStructured(
	ctx context.Context,
	format OutputFormat,
	sysPrompt string,
	question string,
	topic string,
) (*FinalResponse, error) {
	orchStart := time.Now()

	o.log.Info("orchestrator_start",
		zap.String("format", string(format)),
		zap.String("question_preview", truncate(question, 80)),
		zap.String("topic", topic),
	)

	// ── Step 1: Call Qwen text model ──────────────────────────────────────────
	chatResp, err := o.textClient.GenerateChatCompletion(ctx, ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: question},
		},
		MaxTokens:   2048,
		Temperature: 0.5,
	})
	if err != nil {
		return nil, fmt.Errorf("text generation failed: %w", err)
	}

	rawContent := chatResp.Content
	tokens := chatResp.TokensUsed

	o.log.Info("raw_ai_response_received",
		zap.String("format", string(format)),
		zap.Int("tokens", tokens),
		zap.String("raw_preview", truncate(rawContent, 200)),
	)

	// ── Step 2: Parse JSON into typed struct + topic validation + retry ────────
	final, parseSuccess, retried := o.parseWithRetry(ctx, format, sysPrompt, question, rawContent, topic)
	if retried && final != nil {
		tokens += final.TokensUsed // retry consumed more tokens
	}

	o.log.Info("json_parse_result",
		zap.Bool("json_parse_success", parseSuccess),
		zap.Bool("retried", retried),
		zap.String("format", string(format)),
	)

	final.TokensUsed = tokens
	final.IsStructuredJSON = parseSuccess
	final.JSONParseSuccess = parseSuccess
	final.Mode = format

	// ── Step 3: Image generation for visual formats ───────────────────────────
	isVisual := format == OutputFormatAnime || format == OutputFormatSports
	final.IllustrationURL = ""

	if isVisual && o.imageGen != nil {
		scenePrompt := o.extractScenePromptFromFinal(final, format)

		o.log.Info("image_generation_triggered",
			zap.Bool("image_triggered", scenePrompt != ""),
			zap.String("format", string(format)),
			zap.String("image_prompt", truncate(scenePrompt, 120)),
		)

		if scenePrompt != "" {
			style := ImageStyleAnime
			if format == OutputFormatSports {
				style = ImageStyleSports
			}
			fullPrompt := BuildImagePrompt(style, scenePrompt)

			imgCtx, imgCancel := context.WithTimeout(context.Background(), 80*time.Second)
			defer imgCancel()

			imgStart := time.Now()
			imgResult := o.imageGen.GenerateImage(imgCtx, fullPrompt, style)
			imgDuration := int(time.Since(imgStart).Milliseconds())

			if !imgResult.Fallback {
				final.IllustrationURL = strings.TrimSpace(imgResult.ImageURL)
				o.log.Info("image_generation_complete",
					zap.String("image_url", final.IllustrationURL),
					zap.Int("image_generation_duration", imgDuration),
					zap.Bool("fallback_triggered", false),
				)
			} else {
				o.log.Warn("image_generation_fallback",
					zap.Bool("fallback_triggered", true),
					zap.String("image_error", imgResult.FallbackMsg),
					zap.Int("image_generation_duration", imgDuration),
				)
			}
		} else {
			o.log.Warn("image_generation_skipped — no visual_scene_prompt in response",
				zap.String("format", string(format)),
			)
		}
	} else {
		o.log.Debug("image_generation_not_applicable",
			zap.Bool("image_triggered", false),
			zap.String("format", string(format)),
		)
	}

	orchestrationMs := int(time.Since(orchStart).Milliseconds())
	final.OrchestrationMs = orchestrationMs

	o.log.Info("orchestrator_complete",
		zap.String("format", string(format)),
		zap.Bool("json_parse_success", parseSuccess),
		zap.Bool("has_illustration", final.IllustrationURL != ""),
		zap.String("illustration_url", final.IllustrationURL),
		zap.Int("orchestration_time", orchestrationMs),
		zap.Int("tokens_used", final.TokensUsed),
	)

	return final, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Internal: parse with retry
// ──────────────────────────────────────────────────────────────────────────────

// parseWithRetry attempts to parse the AI raw response into typed structs.
// If parsing fails, it retries once with a correction prompt.
// Always returns a non-nil FinalResponse.
func (o *AIOrchestrator) parseWithRetry(
	ctx context.Context,
	format OutputFormat,
	sysPrompt, question, rawContent, topic string,
) (final *FinalResponse, success bool, retried bool) {
	// Try parsing the initial response
	final, success = o.parseRaw(format, rawContent)
	isValidTopic := isTopicValid(rawContent, topic)

	if success && isValidTopic {
		return final, true, false
	}

	// ── Retry once with correction prompt ─────────────────────────────────────
	o.log.Warn("parse_or_validation_failed — retrying",
		zap.String("format", string(format)),
		zap.Bool("json_parse_success", success),
		zap.Bool("topic_valid", isValidTopic),
	)

	var correctionPrompt string
	if !isValidTopic {
		correctionPrompt = fmt.Sprintf(
			"The previous answer drifted from the topic.\n"+
				"Re-answer strictly about %s.\n"+
				"Do not introduce unrelated elements.\n\n"+
				"Required fields for format '%s':\n%s\n\n"+
				"Start your response with { and end with }",
			topic, format, schemaHint(format),
		)
	} else {
		correctionPrompt = fmt.Sprintf(
			"CRITICAL: Your previous response was not valid JSON or was missing required fields.\n\n"+
				"You MUST return ONLY a single raw JSON object. No markdown code fences, no explanation text, no apology.\n"+
				"Required fields for format '%s':\n%s\n\n"+
				"Start your response with { and end with }",
			format, schemaHint(format),
		)
	}

	retryResp, retryErr := o.textClient.GenerateChatCompletion(ctx, ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: question},
			{Role: "assistant", Content: rawContent},
			{Role: "user", Content: correctionPrompt},
		},
		MaxTokens:   2048,
		Temperature: 0.1, // Very low temp for deterministic JSON
	})
	if retryErr != nil {
		o.log.Warn("json_parse_retry_call_failed", zap.Error(retryErr))
		// Return a fallback FinalResponse with the raw text
		return rawFallback(format, rawContent, 0), false, true
	}

	retryTokens := retryResp.TokensUsed
	final, success = o.parseRaw(format, retryResp.Content)
	isValidTopicRetry := isTopicValid(retryResp.Content, topic)

	if success && isValidTopicRetry {
		final.TokensUsed = retryTokens
		o.log.Info("parse_and_validation_succeeded_after_retry")
		return final, true, true
	}

	// Both attempts failed — use safe fallback
	o.log.Warn("validation_failed_after_retry — using safe fallback",
		zap.Bool("json_parse_success", success),
		zap.Bool("topic_valid", isValidTopicRetry),
	)
	fallbackMsg := "Let's refocus on the physics concept."
	fallback := rawFallback(format, fallbackMsg, retryTokens)
	return fallback, false, true
}

func isTopicValid(content, topic string) bool {
	if topic == "" {
		return true
	}
	lowerC := strings.ToLower(content)
	lowerT := strings.ToLower(topic)

	// Check exact match
	if strings.Contains(lowerC, lowerT) {
		return true
	}

	// Check if significant keywords exist
	words := strings.Fields(lowerT)
	for _, w := range words {
		if len(w) > 3 && strings.Contains(lowerC, w) {
			return true
		}
	}
	// Strict: If we totally can't find core keywords, reject.
	return false
}

// parseRaw attempts to unmarshal raw AI output into the correct typed struct.
func (o *AIOrchestrator) parseRaw(format OutputFormat, raw string) (*FinalResponse, bool) {
	// Strip markdown code fences
	s := strings.TrimSpace(raw)
	for _, fence := range []string{"```json", "```"} {
		if strings.HasPrefix(s, fence) {
			s = strings.TrimPrefix(s, fence)
			if idx := strings.LastIndex(s, "```"); idx >= 0 {
				s = s[:idx]
			}
			s = strings.TrimSpace(s)
			break
		}
	}

	// Find JSON object boundary (in case model prepended text)
	if start := strings.Index(s, "{"); start > 0 {
		s = s[start:]
	}

	result := &FinalResponse{RawJSON: s}

	switch format {
	case OutputFormatAnime:
		var v AnimeResponse
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return result, false
		}
		if v.EpisodeTitle == "" || v.PhysicsExplanation == "" {
			return result, false
		}
		result.Anime = &v
		result.Title = v.EpisodeTitle
		result.MainContent = v.PhysicsExplanation
		return result, true

	case OutputFormatSports:
		var v SportsResponse
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return result, false
		}
		if v.GameTitle == "" || v.PlayBreakdown == "" {
			return result, false
		}
		result.Sports = &v
		result.Title = v.GameTitle
		result.MainContent = v.PlayBreakdown
		return result, true

	case OutputFormatSummary:
		var v SummaryResponse
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return result, false
		}
		if v.Title == "" || v.CoreConceptExplanation == "" {
			return result, false
		}
		result.Summary = &v
		result.Title = v.Title
		result.MainContent = v.CoreConceptExplanation
		return result, true

	case OutputFormatDetailed:
		var v DetailedResponse
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return result, false
		}
		if v.Title == "" || v.ConceptExplanation == "" {
			return result, false
		}
		result.Detailed = &v
		result.Title = v.Title
		result.MainContent = v.ConceptExplanation
		return result, true

	case OutputFormatAcademic:
		var v AcademicResponse
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return result, false
		}
		if v.Title == "" {
			return result, false
		}
		result.Academic = &v
		result.Title = v.Title
		result.MainContent = v.TheoreticalBackground
		return result, true
	}

	return result, false
}

// extractScenePromptFromFinal pulls the visual_scene_prompt from the typed struct.
func (o *AIOrchestrator) extractScenePromptFromFinal(f *FinalResponse, format OutputFormat) string {
	switch format {
	case OutputFormatAnime:
		if f.Anime != nil {
			return strings.TrimSpace(f.Anime.VisualScenePrompt)
		}
		// Fallback: try raw JSON
		return ExtractScenePrompt(f.RawJSON)
	case OutputFormatSports:
		if f.Sports != nil {
			return strings.TrimSpace(f.Sports.VisualScenePrompt)
		}
		return ExtractScenePrompt(f.RawJSON)
	}
	return ""
}

// rawFallback creates a FinalResponse for when JSON parsing completely fails.
// The raw text is exposed as MainContent so the user at least sees something.
func rawFallback(format OutputFormat, rawText string, tokens int) *FinalResponse {
	return &FinalResponse{
		Mode:             format,
		RawJSON:          rawText,
		Title:            "",
		MainContent:      rawText,
		IsStructuredJSON: false,
		JSONParseSuccess: false,
		TokensUsed:       tokens,
	}
}

// schemaHint returns a short field list for the correction prompt.
func schemaHint(format OutputFormat) string {
	switch format {
	case OutputFormatAnime:
		return `"episode_title", "main_character", "story_arc", "physics_explanation", "visual_scene_prompt"`
	case OutputFormatSports:
		return `"game_title", "sport_used", "play_breakdown", "coaching_tip", "scoreboard_summary", "visual_scene_prompt"`
	case OutputFormatSummary:
		return `"title", "core_concept_explanation", "key_points" (array), "real_world_example", "short_conclusion"`
	case OutputFormatDetailed:
		return `"title", "concept_explanation", "step_by_step_breakdown" (array), "example", "mini_quiz" (array)`
	case OutputFormatAcademic:
		return `"title", "abstract", "theoretical_background", "methodology", "conclusion"`
	}
	return "valid JSON object"
}
