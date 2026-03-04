package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OutputFormat enumerates every supported output mode.
type OutputFormat string

const (
	OutputFormatSummary     OutputFormat = "quick_summary"
	OutputFormatDetailed    OutputFormat = "detailed_explanation"
	OutputFormatAnime       OutputFormat = "anime_style"
	OutputFormatSports      OutputFormat = "sports_analogy"
	OutputFormatAcademic    OutputFormat = "academic_formal"
	OutputFormatSlides      OutputFormat = "presentation_slides"
	OutputFormatAudio       OutputFormat = "audio_explanation"
	OutputFormatTranslation OutputFormat = "translation"
)

// OutputFormatLabel maps internal keys to human-readable labels.
var OutputFormatLabel = map[OutputFormat]string{
	OutputFormatSummary:     "Quick Summary",
	OutputFormatDetailed:    "Detailed Explanation",
	OutputFormatAnime:       "Anime Style Illustration",
	OutputFormatSports:      "Sports Analogy",
	OutputFormatAcademic:    "Academic Formal",
	OutputFormatSlides:      "Presentation Slides",
	OutputFormatAudio:       "Audio Explanation",
	OutputFormatTranslation: "Translation",
}

// PromptBuilderConfig carries all personalisation signals.
type PromptBuilderConfig struct {
	LearningStyle    string // visual, auditory, reading, kinesthetic, adaptive
	DominantInterest string // anime, sports, academic, music, gaming, …
	ExplanationDepth string // beginner, intermediate, advanced
	OutputFormat     OutputFormat
	Topic            string // locked topic extracted from user query
	Domain           string // locked domain isolated from query
	RetrievedContext string // RAG injected text (already sanitised)
	TargetLanguage   string // for translation mode
}

// ──────────────────────────────────────────────────────────────────────────────
// Strict JSON Schemas per Format
// ──────────────────────────────────────────────────────────────────────────────

const schemaPrefix = `CRITICAL: You MUST return ONLY a single valid JSON object with NO markdown, NO code fences, NO commentary, NO apology, NO preamble. Output only raw JSON.

`

const schemaSummary = schemaPrefix + `Schema:
{
  "title": "<string: concise concept title>",
  "core_concept_explanation": "<string: clear comprehensive paragraph explaining the core idea>",
  "key_points": ["<string>", "..."],
  "real_world_example": "<string: structured real-world scenario or an example>",
  "short_conclusion": "<string: brief concluding thought or recap>",
  "visual_scene_prompt": "<string: vivid 1 sentence scene for image generation, no special chars>"
}
Rules: Provide an educational depth. DO NOT make it overly compressed. core_concept_explanation MUST be a complete paragraph. key_points MUST have 3-5 entries, all strings. visual_scene_prompt MUST be a single descriptive sentence suitable for image generation.`

const schemaDetailed = schemaPrefix + `Schema:
{
  "title": "<string: concept title>",
  "concept_explanation": "<string: clear multi-paragraph explanation>",
  "step_by_step_breakdown": ["<step 1>", "<step 2>", "..."],
  "example": "<string: concrete worked example>",
  "mini_quiz": ["<question 1>?", "<question 2>?", "<question 3>?"],
  "visual_scene_prompt": "<string: vivid 1 sentence scene for image generation, no special chars>"
}
Rules: step_by_step_breakdown MUST have 3-7 entries. mini_quiz MUST have exactly 3 questions. visual_scene_prompt MUST be a single descriptive sentence suitable for image generation.`

const schemaAnime = schemaPrefix + `Schema:
{
  "episode_title": "<string: dramatic anime episode title>",
  "main_character": "<string: protagonist name and brief description>",
  "story_arc": "<string: 2-3 sentence narrative connecting concept to anime story>",
  "physics_explanation": "<string: accurate concept explanation woven into story>",
  "visual_scene_prompt": "<string: vivid 1 sentence scene for image generation, no special chars>"
}
Rules: Explain the physics concept using short anime-themed analogy. Do not create unrelated narrative arcs. All story elements must directly map to the physics principle. visual_scene_prompt MUST be a single descriptive sentence suitable for image generation.`

const schemaSports = schemaPrefix + `Schema:
{
  "game_title": "<string: sports scenario title>",
  "sport_used": "<string: which sport is the analogy>",
  "play_breakdown": "<string: the concept explained as a sports play>",
  "coaching_tip": "<string: key insight framed as coach advice>",
  "scoreboard_summary": "<string: key formula or rule as a scoreboard stat>",
  "visual_scene_prompt": "<string: vivid 1 sentence sports scene for image generation, no special chars>"
}
Rules: visual_scene_prompt MUST be a single descriptive sentence suitable for image generation.`

const schemaAcademic = schemaPrefix + `Schema:
{
  "title": "<string: formal academic title>",
  "abstract": "<string: 2-3 sentence abstract>",
  "theoretical_background": "<string: formal explanation with definitions>",
  "methodology": "<string: step-by-step formal breakdown>",
  "conclusion": "<string: formal summary and implications>",
  "visual_scene_prompt": "<string: vivid 1 sentence scene for image generation, no special chars>"
}
Rules: visual_scene_prompt MUST be a single descriptive sentence suitable for image generation.`

// BuildSystemPrompt returns the full system prompt string.
func BuildSystemPrompt(cfg PromptBuilderConfig) string {
	var sb strings.Builder

	sb.WriteString("You are an educational physics tutor.\n")
	sb.WriteString("Explain physics concepts clearly.\n")
	sb.WriteString("If the user asks for an illustration, generate a visual_scene_prompt describing a clear visual scene that explains the physics concept.\n")
	sb.WriteString("Stay on topic and do not introduce unrelated concepts.\n\n")

	// RAG context injection (BEFORE schema so model has context when writing JSON)
	if strings.TrimSpace(cfg.RetrievedContext) != "" {
		sb.WriteString("## Context from user's uploaded materials:\n")
		sb.WriteString("---\n")
		sb.WriteString(cfg.RetrievedContext)
		sb.WriteString("\n---\n")
		sb.WriteString("Use the above context when forming your answer. If the answer is not in the context, say so clearly.\n\n")
	}

	// Strict schema per output format
	switch cfg.OutputFormat {
	case OutputFormatSummary:
		sb.WriteString(schemaSummary)
	case OutputFormatDetailed:
		sb.WriteString(schemaDetailed)
	case OutputFormatAnime:
		sb.WriteString(schemaAnime)
	case OutputFormatSports:
		sb.WriteString(schemaSports)
	case OutputFormatAcademic:
		sb.WriteString(schemaAcademic)
	case OutputFormatSlides:
		sb.WriteString(schemaPrefix + `Schema:
{"slides":[{"title":"...","content":"...","speaker_notes":"..."}]}
Rules: 5-8 slides. title ≤ 60 chars. content ≤ 200 chars bullet points. speaker_notes ≤ 300 chars.`)
	case OutputFormatAudio:
		sb.WriteString("OUTPUT FORMAT: Write a natural spoken-word script only. Short sentences. No markdown. Conversational rhythm. Target: 2-3 minutes of speech. No JSON.\n")
	case OutputFormatTranslation:
		lang := coalesce(cfg.TargetLanguage, "Indonesian")
		sb.WriteString(fmt.Sprintf("Translate the following into %s. Preserve structure and formatting. Return ONLY the translated text, no commentary.\n", lang))
	default:
		sb.WriteString("Respond clearly and helpfully. No markdown unless necessary.\n")
	}

	return sb.String()
}

// BuildOutputFormatPromptRequest returns the canned message asking the user to
// choose their preferred output mode.
func BuildOutputFormatPromptRequest() string {
	return `How would you like the material delivered? Please choose:

1. 📝 Quick Summary
2. 📖 Detailed Explanation
3. 🎌 Anime Style Illustration
4. ⚽ Sports Analogy
5. 🎓 Academic Formal
6. 📊 Presentation Slides
7. 🔊 Audio Explanation
8. 🌍 Translation

Reply with the number or name of your preferred format.`
}

// ParseOutputFormat converts a user reply (number or name) into an OutputFormat.
func ParseOutputFormat(input string) (OutputFormat, bool) {
	switch strings.TrimSpace(strings.ToLower(input)) {
	case "1", "quick summary", "summary":
		return OutputFormatSummary, true
	case "2", "detailed explanation", "detailed":
		return OutputFormatDetailed, true
	case "3", "anime style illustration", "anime":
		return OutputFormatAnime, true
	case "4", "sports analogy", "sports":
		return OutputFormatSports, true
	case "5", "academic formal", "academic":
		return OutputFormatAcademic, true
	case "6", "presentation slides", "slides", "ppt":
		return OutputFormatSlides, true
	case "7", "audio explanation", "audio":
		return OutputFormatAudio, true
	case "8", "translation":
		return OutputFormatTranslation, true
	}
	return "", false
}

// ──────────────────────────────────────────────────────────────────────────────
// JSON Schema Validation Layer
// ──────────────────────────────────────────────────────────────────────────────

// ValidateStructuredJSON validates that the AI response conforms to the
// expected JSON schema for the given OutputFormat. Returns the sanitised
// JSON string or an error if validation fails.
func ValidateStructuredJSON(raw string, format OutputFormat) (string, error) {
	// Strip possible markdown code fence
	s := strings.TrimSpace(raw)
	for _, fence := range []string{"```json", "```"} {
		if strings.HasPrefix(s, fence) {
			s = s[len(fence):]
			if idx := strings.LastIndex(s, "```"); idx >= 0 {
				s = s[:idx]
			}
			s = strings.TrimSpace(s)
			break
		}
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(s), &parsed); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	// Format-specific field checks
	switch format {
	case OutputFormatSummary:
		if err := requireFields(parsed, "title", "core_concept_explanation", "key_points", "real_world_example", "short_conclusion", "visual_scene_prompt"); err != nil {
			return "", err
		}
	case OutputFormatDetailed:
		if err := requireFields(parsed, "title", "concept_explanation", "step_by_step_breakdown", "example", "mini_quiz", "visual_scene_prompt"); err != nil {
			return "", err
		}
	case OutputFormatAnime:
		if err := requireFields(parsed, "episode_title", "main_character", "story_arc", "physics_explanation", "visual_scene_prompt"); err != nil {
			return "", err
		}
	case OutputFormatSports:
		if err := requireFields(parsed, "game_title", "sport_used", "play_breakdown", "coaching_tip", "scoreboard_summary", "visual_scene_prompt"); err != nil {
			return "", err
		}
	case OutputFormatAcademic:
		if err := requireFields(parsed, "title", "abstract", "theoretical_background", "methodology", "conclusion", "visual_scene_prompt"); err != nil {
			return "", err
		}
	}

	// Canonicalise back to compact JSON
	out, err := json.Marshal(parsed)
	if err != nil {
		return "", fmt.Errorf("re-serialising JSON: %w", err)
	}
	return string(out), nil
}

// ExtractScenePrompt extracts the visual_scene_prompt field from an anime or sports JSON blob.
func ExtractScenePrompt(raw string) string {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return ""
	}
	if v, ok := parsed["visual_scene_prompt"].(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

// requireFields returns an error if any field key is missing.
func requireFields(m map[string]interface{}, fields ...string) error {
	for _, f := range fields {
		if _, ok := m[f]; !ok {
			return fmt.Errorf("missing required field: %q", f)
		}
	}
	return nil
}

func coalesce(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
