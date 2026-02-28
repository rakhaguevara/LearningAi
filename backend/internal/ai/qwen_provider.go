package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// QwenProvider implements AIProvider using Alibaba Cloud's Qwen LLM.
// Currently uses mock responses; swap internals for real API calls in production.
type QwenProvider struct {
	apiKey   string
	endpoint string
	model    string
	log      *zap.Logger
}

func NewQwenProvider(apiKey, endpoint, model string, log *zap.Logger) *QwenProvider {
	return &QwenProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
		model:    model,
		log:      log,
	}
}

func (q *QwenProvider) ExplainConcept(ctx context.Context, req ExplainRequest) (*ExplainResponse, error) {
	start := time.Now()

	systemPrompt := q.buildSystemPrompt(req.Interests, req.Style, req.Difficulty)

	q.log.Info("explaining concept",
		zap.String("topic", req.Topic),
		zap.String("subject", req.Subject),
		zap.String("model", q.model),
	)

	// TODO: Replace with real Qwen API call
	// Real implementation would POST to q.endpoint with:
	// - model: q.model
	// - messages: [{role: "system", content: systemPrompt}, {role: "user", content: userPrompt}]
	// - Authorization: Bearer q.apiKey

	// userPrompt will be sent to the Qwen API in production
	_ = fmt.Sprintf(
		"Explain the concept of '%s' in the subject of '%s'. "+
			"Make it engaging and use analogies from these interests: %s",
		req.Topic, req.Subject, strings.Join(req.Interests, ", "),
	)

	interestRef := "general"
	if len(req.Interests) > 0 {
		interestRef = req.Interests[0]
	}

	explanation := fmt.Sprintf(
		"[Mock Qwen Response]\n\n"+
			"System: %s\n\n"+
			"Topic: %s | Subject: %s\n\n"+
			"Imagine you're in a world where %s meets %s. "+
			"The concept of %s works like a power-up system in your favorite game — "+
			"each level builds on the last, unlocking new abilities and understanding. "+
			"Just like how a protagonist trains to master new skills, learning %s "+
			"requires building foundational knowledge first, then layering complexity.\n\n"+
			"This is a mock response. Connect the Qwen API to generate real adaptive explanations.",
		systemPrompt,
		req.Topic, req.Subject,
		interestRef, req.Subject,
		req.Topic, req.Topic,
	)

	latency := time.Since(start).Milliseconds()

	return &ExplainResponse{
		Explanation: explanation,
		TokensUsed:  150,
		LatencyMs:   int(latency),
	}, nil
}

func (q *QwenProvider) GenerateIllustration(ctx context.Context, req IllustrationRequest) (*IllustrationResponse, error) {
	q.log.Info("generating illustration",
		zap.String("topic", req.Topic),
		zap.String("style", req.Style),
	)

	// TODO: Replace with real Qwen VL / image generation API call
	prompt := fmt.Sprintf(
		"Create an educational illustration about '%s'. "+
			"Style: %s. Incorporate visual themes from: %s. "+
			"Description: %s. "+
			"The image should be clear, educational, and visually engaging.",
		req.Topic,
		req.Style,
		strings.Join(req.Interests, ", "),
		req.Description,
	)

	return &IllustrationResponse{
		ImageURL:   fmt.Sprintf("https://placeholder.ailearn.dev/illustrations/%s.png", strings.ReplaceAll(req.Topic, " ", "-")),
		Prompt:     prompt,
		TokensUsed: 80,
	}, nil
}

func (q *QwenProvider) AdaptTeachingStyle(ctx context.Context, req StyleRequest) (*StyleResponse, error) {
	q.log.Info("adapting teaching style",
		zap.Strings("interests", req.Interests),
		zap.String("difficulty", req.DifficultyLevel),
	)

	interestContext := strings.Join(req.Interests, ", ")

	systemPrompt := fmt.Sprintf(
		"You are an adaptive AI tutor. The student is interested in: %s. "+
			"Their difficulty level is: %s. Their preferred learning style is: %s. "+
			"Use analogies, metaphors, and examples drawn from their interests to explain concepts. "+
			"Keep the tone encouraging and match the complexity to their level. "+
			"If they like anime, use anime references. If they like sports, use sports metaphors. "+
			"Always make learning feel relevant and exciting.",
		interestContext, req.DifficultyLevel, req.PreferredStyle,
	)

	tone := "encouraging and relatable"
	if req.DifficultyLevel == "advanced" {
		tone = "precise and technically deep"
	} else if req.DifficultyLevel == "beginner" {
		tone = "simple, warm, and reassuring"
	}

	return &StyleResponse{
		SystemPrompt: systemPrompt,
		Tone:         tone,
		Examples:     fmt.Sprintf("Use references from: %s", interestContext),
	}, nil
}

func (q *QwenProvider) buildSystemPrompt(interests []string, style, difficulty string) string {
	if len(interests) == 0 {
		interests = []string{"general knowledge"}
	}
	return fmt.Sprintf(
		"You are an adaptive learning assistant. Tailor explanations using themes from: %s. "+
			"Style: %s. Difficulty: %s. Be engaging, accurate, and creative.",
		strings.Join(interests, ", "), style, difficulty,
	)
}
