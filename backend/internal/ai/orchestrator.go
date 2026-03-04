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
	imageGen   *AIImageService
	log        *zap.Logger
}

// NewAIOrchestrator creates the orchestrator.
func NewAIOrchestrator(textClient *QwenClient, imageGen *AIImageService, log *zap.Logger) *AIOrchestrator {
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
	history []ChatMessage,
) (*FinalResponse, error) {
	orchStart := time.Now()

	o.log.Info("orchestrator_start",
		zap.String("format", string(format)),
		zap.String("question_preview", truncate(question, 80)),
	)

	messages := make([]ChatMessage, 0, len(history)+2)
	messages = append(messages, ChatMessage{Role: "system", Content: sysPrompt})
	messages = append(messages, history...)
	messages = append(messages, ChatMessage{Role: "user", Content: question})

	resp, err := o.textClient.GenerateChatCompletion(ctx, ChatRequest{
		Messages:     messages,
		MaxTokens:    1500,
		Temperature:  0.2,  // Prevents hallucination
		ResponseJSON: true, // Native JSON output
	})
	if err != nil {
		return nil, fmt.Errorf("text generation failed: %w", err)
	}

	rawContent := resp.Content
	tokens := resp.TokensUsed

	o.log.Info("raw_ai_response_received",
		zap.String("format", string(format)),
		zap.Int("tokens", tokens),
		zap.String("raw_preview", truncate(rawContent, 200)),
	)

	// Parse JSON directly (no retry loop needed since we forced native JSON)
	final, parseSuccess := o.parseRaw(format, rawContent)

	o.log.Info("json_parse_result",
		zap.Bool("json_parse_success", parseSuccess),
		zap.String("format", string(format)),
	)

	if !parseSuccess || final == nil {
		// Native json mode failed us (rare)
		final = rawFallback(format, rawContent, tokens)
	}

	final.TokensUsed = tokens
	final.IsStructuredJSON = parseSuccess
	final.JSONParseSuccess = parseSuccess
	final.Mode = format

	// ── Step 3: Image generation ──────────────────────────────────────────────
	final.IllustrationURL = ""

	// Extract all image generation fields from the tutor JSON response
	scenePrompt := ""
	imageStyle := ""
	negativePrompt := ""
	var diagramLabels []string

	var rawMap map[string]interface{}
	if err := json.Unmarshal([]byte(final.RawJSON), &rawMap); err == nil {
		if v, ok := rawMap["visual_scene_prompt"].(string); ok {
			scenePrompt = strings.TrimSpace(v)
		}
		if v, ok := rawMap["image_style"].(string); ok {
			imageStyle = strings.TrimSpace(v)
		}
		if v, ok := rawMap["negative_prompt"].(string); ok {
			negativePrompt = strings.TrimSpace(v)
		}
		if v, ok := rawMap["diagram_labels"].([]interface{}); ok {
			for _, lbl := range v {
				if s, ok := lbl.(string); ok && s != "" {
					diagramLabels = append(diagramLabels, s)
				}
			}
		}
	}

	o.log.Info("image_generation_fields_extracted",
		zap.String("scene_prompt_preview", truncate(scenePrompt, 100)),
		zap.String("image_style", imageStyle),
		zap.String("negative_prompt_preview", truncate(negativePrompt, 50)),
		zap.Int("diagram_labels_count", len(diagramLabels)),
	)

	if scenePrompt == "" {
		o.log.Warn("image_generation_fallback_triggered",
			zap.String("reason", "visual_scene_prompt missing from generated JSON"),
			zap.String("topic_fallback", topic),
		)
		// Fallback: construct minimal scene prompt from topic
		scenePrompt = fmt.Sprintf("educational illustration showing %s with clear labeled components", coalesce(topic, "the concept previously discussed"))
	}

	if o.imageGen != nil {
		// Build fully enhanced image generation input via the prompt enhancer
		imgInput := BuildImageGenerationInput(scenePrompt, imageStyle, negativePrompt, diagramLabels, topic)

		o.log.Info("image_generation_triggered",
			zap.String("format", string(format)),
			zap.String("style", imgInput.Style),
			zap.String("final_prompt_preview", truncate(imgInput.FinalPrompt, 150)),
			zap.String("negative_prompt_preview", truncate(imgInput.NegativePrompt, 80)),
		)

		imgCtx, imgCancel := context.WithTimeout(context.Background(), 80*time.Second)
		defer imgCancel()

		imgStart := time.Now()
		imgResult, err := o.imageGen.GenerateImageFromInput(imgCtx, imgInput)
		imgDuration := int(time.Since(imgStart).Milliseconds())

		// Log detailed result information
		if imgResult != nil {
			o.log.Info("image_generation_result_received",
				zap.Bool("has_result", true),
				zap.Bool("result_fallback", imgResult.Fallback),
				zap.String("fallback_msg", imgResult.FallbackMsg),
				zap.String("image_url", imgResult.ImageURL),
				zap.Int("duration_ms", imgResult.DurationMs),
			)
		} else {
			o.log.Warn("image_generation_result_nil",
				zap.Bool("has_result", false),
				zap.Error(err),
			)
		}

		if err == nil && imgResult != nil && !imgResult.Fallback {
			final.IllustrationURL = strings.TrimSpace(imgResult.ImageURL)
			o.log.Info("image_generation_complete",
				zap.String("image_url", final.IllustrationURL),
				zap.Int("image_generation_duration", imgDuration),
				zap.Bool("fallback_triggered", false),
			)
		} else {
			errMsg := "unknown error"
			if err != nil {
				errMsg = err.Error()
			} else if imgResult != nil {
				errMsg = imgResult.FallbackMsg
			}
			o.log.Warn("image_generation_fallback",
				zap.Bool("fallback_triggered", true),
				zap.String("image_error", errMsg),
				zap.Int("image_generation_duration", imgDuration),
				zap.String("error_summary", fmt.Sprintf("Image gen failed: %s", errMsg)),
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

// (Legacy Retry functions removed since Native JSON mode guarantees structure)

// parseRaw attempts to unmarshal raw AI output into the correct typed struct.
func (o *AIOrchestrator) parseRaw(format OutputFormat, raw string) (*FinalResponse, bool) {
	// Native JSON mode guarantees JSON format, no regex stripping needed.
	s := strings.TrimSpace(raw)

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
		result.IllustrationURL = v.VisualScenePrompt // mapped internally, but will be overwritten by the real generated OSS URL later
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
		result.IllustrationURL = v.VisualScenePrompt
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

// (Old extractScenePromptFromFinal removed in favor of direct json map lookup)

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
