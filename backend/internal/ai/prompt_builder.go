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
// Image generation constants
// ──────────────────────────────────────────────────────────────────────────────

// DefaultImageStyle is the consistent illustration style applied to every image.
const DefaultImageStyle = "modern flat educational illustration, soft lighting, white or light gradient background, vector-like clarity, textbook style, clean and professional"

// DefaultNegativePrompt is always appended to prevent bad outputs.
const DefaultNegativePrompt = "blurry, distorted, extra limbs, messy background, dark lighting, low resolution, abstract art, photorealistic, watermark, text overlay, noise, grainy"

// MinVisualPromptWords is the threshold below which the prompt enhancer fires.
const MinVisualPromptWords = 40

// ──────────────────────────────────────────────────────────────────────────────
// Strict JSON Schemas per Format
// ──────────────────────────────────────────────────────────────────────────────

const schemaPrefix = `CRITICAL: You MUST return ONLY a single valid JSON object with NO markdown, NO code fences, NO commentary, NO apology, NO preamble. Output only raw JSON.

`

// imageFieldsSchema is appended to every schema to enforce consistent image generation fields.
const imageFieldsSchema = `
  "visual_scene_prompt": "<string: HIGHLY DESCRIPTIVE educational illustration scene. Minimum 40 words. Describe: objects, their positions relative to each other, lighting direction (soft diffused light from upper-left), camera angle (eye-level or slight top-down for diagrams), background (clean white or light gradient), color palette (bright and clear), and how the scene visually communicates the physics concept. Avoid abstract art. Focus on educational clarity.>",
  "image_style": "<string: one of — 'flat_educational', 'diagram_technical', 'infographic_clean', 'textbook_illustration'>",
  "negative_prompt": "<string: comma-separated list of things to avoid, always include: blurry, distorted, dark lighting, messy background, abstract art, photorealistic>",
  "diagram_labels": ["<string: label for a key element in the scene>", "..."]`

const schemaSummary = schemaPrefix + `Schema:
{
  "title": "<string: concise concept title>",
  "core_concept_explanation": "<string: clear comprehensive paragraph explaining the core idea>",
  "key_points": ["<string>", "..."],
  "real_world_example": "<string: structured real-world scenario or an example>",
  "short_conclusion": "<string: brief concluding thought or recap>",` +
	imageFieldsSchema + `
}
Rules: core_concept_explanation MUST be a complete paragraph. key_points MUST have 3-5 entries. visual_scene_prompt MUST be at least 40 words and describe a clear educational scene with spatial layout, lighting, and how it represents the concept. diagram_labels MUST have 2-5 entries naming key visual elements.`

const schemaDetailed = schemaPrefix + `Schema:
{
  "title": "<string: concept title>",
  "concept_explanation": "<string: clear multi-paragraph explanation>",
  "step_by_step_breakdown": ["<step 1>", "<step 2>", "..."],
  "example": "<string: concrete worked example>",
  "mini_quiz": ["<question 1>?", "<question 2>?", "<question 3>?"],` +
	imageFieldsSchema + `
}
Rules: step_by_step_breakdown MUST have 3-7 entries. mini_quiz MUST have exactly 3 questions. visual_scene_prompt MUST be at least 40 words describing a clear educational scene. diagram_labels MUST have 2-5 entries.`

const schemaAnime = schemaPrefix + `Schema:
{
  "episode_title": "<string: dramatic anime episode title>",
  "main_character": "<string: protagonist name and brief description>",
  "story_arc": "<string: 2-3 sentence narrative connecting concept to anime story>",
  "physics_explanation": "<string: accurate concept explanation woven into story>",` +
	imageFieldsSchema + `
}
Rules: All story elements must directly map to the physics principle. visual_scene_prompt MUST be at least 40 words. image_style MUST be 'flat_educational'. diagram_labels name key scene elements.`

const schemaSports = schemaPrefix + `Schema:
{
  "game_title": "<string: sports scenario title>",
  "sport_used": "<string: which sport is the analogy>",
  "play_breakdown": "<string: the concept explained as a sports play>",
  "coaching_tip": "<string: key insight framed as coach advice>",
  "scoreboard_summary": "<string: key formula or rule as a scoreboard stat>",` +
	imageFieldsSchema + `
}
Rules: visual_scene_prompt MUST be at least 40 words depicting the sports analogy scene clearly. diagram_labels name athletes, forces, or equipment visible in the scene.`

const schemaAcademic = schemaPrefix + `Schema:
{
  "title": "<string: formal academic title>",
  "abstract": "<string: 2-3 sentence abstract>",
  "theoretical_background": "<string: formal explanation with definitions>",
  "methodology": "<string: step-by-step formal breakdown>",
  "conclusion": "<string: formal summary and implications>",` +
	imageFieldsSchema + `
}
Rules: visual_scene_prompt MUST be at least 40 words, describing a precise technical diagram. image_style MUST be 'diagram_technical'. diagram_labels MUST reference specific labeled diagram components.`

// BuildSystemPrompt returns the full system prompt string.
func BuildSystemPrompt(cfg PromptBuilderConfig) string {
	var sb strings.Builder

	sb.WriteString("You are an expert educational tutor specialising in clear, engaging explanations.\n")
	sb.WriteString("Explain concepts accurately and adapt to the learner's level.\n")
	sb.WriteString("Stay on topic and do not introduce unrelated concepts.\n\n")
	sb.WriteString("## Image Generation Rules (MANDATORY)\n")
	sb.WriteString("Every response that includes a visual_scene_prompt MUST follow these rules:\n")
	sb.WriteString("- visual_scene_prompt: minimum 40 words, describes objects, positions, lighting (soft from upper-left), camera angle, clean background\n")
	sb.WriteString("- image_style: choose from 'flat_educational', 'diagram_technical', 'infographic_clean', 'textbook_illustration'\n")
	sb.WriteString("- negative_prompt: always include 'blurry, distorted, dark lighting, messy background, abstract art, photorealistic'\n")
	sb.WriteString("- diagram_labels: 2-5 short labels naming key visual elements in the scene\n")
	sb.WriteString("- DO NOT mix explanation text into visual_scene_prompt — it must be a pure image description\n\n")

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
		if err := requireFields(parsed, "title", "core_concept_explanation", "key_points", "real_world_example", "short_conclusion", "visual_scene_prompt", "image_style", "negative_prompt", "diagram_labels"); err != nil {
			return "", err
		}
	case OutputFormatDetailed:
		if err := requireFields(parsed, "title", "concept_explanation", "step_by_step_breakdown", "example", "mini_quiz", "visual_scene_prompt", "image_style", "negative_prompt", "diagram_labels"); err != nil {
			return "", err
		}
	case OutputFormatAnime:
		if err := requireFields(parsed, "episode_title", "main_character", "story_arc", "physics_explanation", "visual_scene_prompt", "image_style", "negative_prompt", "diagram_labels"); err != nil {
			return "", err
		}
	case OutputFormatSports:
		if err := requireFields(parsed, "game_title", "sport_used", "play_breakdown", "coaching_tip", "scoreboard_summary", "visual_scene_prompt", "image_style", "negative_prompt", "diagram_labels"); err != nil {
			return "", err
		}
	case OutputFormatAcademic:
		if err := requireFields(parsed, "title", "abstract", "theoretical_background", "methodology", "conclusion", "visual_scene_prompt", "image_style", "negative_prompt", "diagram_labels"); err != nil {
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

// ──────────────────────────────────────────────────────────────────────────────
// Prompt Enhancer — auto-expands short visual_scene_prompts
// ──────────────────────────────────────────────────────────────────────────────

// ImageGenerationInput is the fully prepared input for the image model.
// It is produced by BuildImageGenerationInput and consumed by AIImageService.
type ImageGenerationInput struct {
	// FinalPrompt is the enhanced, style-prefixed prompt sent to the image API.
	FinalPrompt string
	// NegativePrompt is the negative guidance sent to the image API.
	NegativePrompt string
	// Style is the resolved illustration style identifier.
	Style string
	// DiagramLabels are label hints extracted from the tutor response.
	DiagramLabels []string
}

// BuildImageGenerationInput takes raw fields from the tutor JSON response and
// produces a complete, enhanced ImageGenerationInput ready for the image model.
//
// Pipeline:
//  1. Validate / fallback image_style
//  2. Auto-enhance visual_scene_prompt if < MinVisualPromptWords words
//  3. Prepend consistent style prefix
//  4. Merge negative_prompt with DefaultNegativePrompt
func BuildImageGenerationInput(
	scenePrompt string,
	imageStyle string,
	negativePrompt string,
	diagramLabels []string,
	topic string,
) ImageGenerationInput {
	// 1. Resolve style with fallback
	style := resolveImageStyle(imageStyle)

	// 2. Auto-enhance if scene prompt is too short
	enhanced := autoEnhancePrompt(scenePrompt, topic, style)

	// 3. Prepend style prefix for consistent visual output
	stylePrefix := imageStyleToPromptPrefix(style)
	finalPrompt := stylePrefix + enhanced

	// Append diagram label hints if present
	if len(diagramLabels) > 0 {
		finalPrompt += fmt.Sprintf(". Key elements labeled: %s", strings.Join(diagramLabels, ", "))
	}

	// 4. Merge negative prompts (deduplicated)
	finalNeg := mergeNegativePrompts(negativePrompt, DefaultNegativePrompt)

	return ImageGenerationInput{
		FinalPrompt:    finalPrompt,
		NegativePrompt: finalNeg,
		Style:          style,
		DiagramLabels:  diagramLabels,
	}
}

// resolveImageStyle validates and normalises the image style string.
func resolveImageStyle(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "flat_educational", "flat educational":
		return "flat_educational"
	case "diagram_technical", "diagram technical":
		return "diagram_technical"
	case "infographic_clean", "infographic clean":
		return "infographic_clean"
	case "textbook_illustration", "textbook illustration":
		return "textbook_illustration"
	default:
		return "flat_educational" // safe default
	}
}

// imageStyleToPromptPrefix maps style identifiers to natural language prefixes
// that guide the image model toward consistent visual outputs.
func imageStyleToPromptPrefix(style string) string {
	switch style {
	case "diagram_technical":
		return "Clean technical diagram, precise line art, white background, labeled components, educational textbook style, vector illustration. "
	case "infographic_clean":
		return "Modern clean infographic, bold colors, flat design, white background, icons and labels, educational poster style. "
	case "textbook_illustration":
		return "Classic textbook illustration, soft pastel colors, white background, cross-section view, detailed labels, educational reference style. "
	default: // flat_educational
		return "Modern flat educational illustration, soft diffused lighting, clean white background, bright clear colors, vector-like quality, textbook clarity. "
	}
}

// autoEnhancePrompt expands the visual_scene_prompt if it is too short.
// Enhancement adds: spatial arrangement, color hints, perspective, and clarity instructions.
func autoEnhancePrompt(prompt, topic, style string) string {
	wordCount := len(strings.Fields(prompt))
	if wordCount >= MinVisualPromptWords {
		return prompt
	}

	// Build enhancement based on topic and style
	topicHint := coalesce(topic, "the educational concept")

	var expansion string
	switch style {
	case "diagram_technical":
		expansion = fmt.Sprintf(
			"The scene shows a precise technical diagram illustrating %s. "+
				"Components are arranged logically left-to-right with clear spacing. "+
				"Arrows indicate direction of forces or flow. "+
				"Each element is distinct with clean outlines. "+
				"Viewed from directly above or straight-on for maximum clarity. "+
				"White background with light gray grid lines for reference. "+
				"Color coding distinguishes different components.",
			topicHint,
		)
	case "infographic_clean":
		expansion = fmt.Sprintf(
			"A clean infographic layout explaining %s. "+
				"Bold geometric shapes frame key information. "+
				"Bright accent colors (blue, green, orange) highlight important elements. "+
				"Icons and simple illustrations complement each point. "+
				"Text labels are minimal and positioned near relevant visuals. "+
				"Viewed from straight-on, white background, ample whitespace.",
			topicHint,
		)
	default: // flat_educational + textbook
		expansion = fmt.Sprintf(
			"An educational illustration showing %s. "+
				"The main subject is centered in the frame with supporting elements arranged symmetrically. "+
				"Soft diffused light comes from the upper-left, creating gentle shadows. "+
				"Objects are depicted with clean flat colors and thin outlines. "+
				"The background is white or very light gradient (pale blue or cream). "+
				"Camera angle is eye-level or slight top-down perspective. "+
				"Every element clearly contributes to understanding the concept.",
			topicHint,
		)
	}

	// Combine original prompt with expansion if original is not empty
	if strings.TrimSpace(prompt) == "" {
		return expansion
	}
	return prompt + ". " + expansion
}

// mergeNegativePrompts combines two negative prompt strings, deduplicating tokens.
func mergeNegativePrompts(custom, base string) string {
	seen := make(map[string]bool)
	var parts []string

	for _, raw := range strings.Split(base+","+custom, ",") {
		token := strings.ToLower(strings.TrimSpace(raw))
		if token != "" && !seen[token] {
			seen[token] = true
			parts = append(parts, strings.TrimSpace(raw))
		}
	}
	return strings.Join(parts, ", ")
}
