package ai

import "context"

// AIProvider defines the abstraction layer for AI services.
// Implementations can target Qwen, OpenAI, or any other LLM backend.
type AIProvider interface {
	// ExplainConcept generates an explanation of a topic adapted to user interests.
	ExplainConcept(ctx context.Context, req ExplainRequest) (*ExplainResponse, error)

	// GenerateIllustration creates a contextual illustration prompt or image URL.
	GenerateIllustration(ctx context.Context, req IllustrationRequest) (*IllustrationResponse, error)

	// AdaptTeachingStyle returns a system prompt tailored to user preferences.
	AdaptTeachingStyle(ctx context.Context, req StyleRequest) (*StyleResponse, error)
}

type ExplainRequest struct {
	Topic      string   `json:"topic" binding:"required"`
	Subject    string   `json:"subject" binding:"required"`
	Interests  []string `json:"interests"`
	Style      string   `json:"style"`
	Difficulty string   `json:"difficulty"`
	UserID     string   `json:"user_id"`
}

type ExplainResponse struct {
	Explanation string `json:"explanation"`
	TokensUsed  int    `json:"tokens_used"`
	LatencyMs   int    `json:"latency_ms"`
}

type IllustrationRequest struct {
	Topic       string   `json:"topic" binding:"required"`
	Description string   `json:"description"`
	Interests   []string `json:"interests"`
	Style       string   `json:"style"`
}

type IllustrationResponse struct {
	ImageURL   string `json:"image_url"`
	Prompt     string `json:"prompt"`
	TokensUsed int    `json:"tokens_used"`
}

type StyleRequest struct {
	Interests       []string `json:"interests"`
	DifficultyLevel string   `json:"difficulty_level"`
	PreferredStyle  string   `json:"preferred_style"`
}

type StyleResponse struct {
	SystemPrompt string `json:"system_prompt"`
	Tone         string `json:"tone"`
	Examples     string `json:"examples"`
}
