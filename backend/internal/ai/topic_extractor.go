package ai

import (
	"context"
	"encoding/json"
	"strings"

	"go.uber.org/zap"
)

// TopicExtractor is responsible for pulling out the core topic and domain
// from a raw user question to enable strict context locking.
type TopicExtractor struct {
	client *QwenClient
	log    *zap.Logger
}

func NewTopicExtractor(client *QwenClient, log *zap.Logger) *TopicExtractor {
	return &TopicExtractor{
		client: client,
		log:    log,
	}
}

// Extract extracts the topic and domain. Returns the raw question if parsing fails.
func (e *TopicExtractor) Extract(ctx context.Context, question string) (topic string, domain string) {
	prompt := `Extract the core educational topic and the domain from this query.
Return ONLY valid JSON matching this schema: {"topic": "<string>", "domain": "<string>"}
Do not include markdown or explanations.`

	resp, err := e.client.GenerateChatCompletion(ctx, ChatRequest{
		Messages: []ChatMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: question},
		},
		MaxTokens:   150,
		Temperature: 0.1,
	})

	topic = question
	domain = "General"

	if err != nil {
		e.log.Warn("topic extraction fell back to query", zap.Error(err))
		return
	}

	content := strings.TrimSpace(resp.Content)
	for _, fence := range []string{"```json", "```"} {
		if strings.HasPrefix(content, fence) {
			content = strings.TrimPrefix(content, fence)
			if idx := strings.LastIndex(content, "```"); idx >= 0 {
				content = content[:idx]
			}
			content = strings.TrimSpace(content)
			break
		}
	}

	if start := strings.Index(content, "{"); start >= 0 {
		content = content[start:]
	}

	var parsed struct {
		Topic  string `json:"topic"`
		Domain string `json:"domain"`
	}

	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		if parsed.Topic != "" {
			topic = parsed.Topic
		}
		if parsed.Domain != "" {
			domain = parsed.Domain
		}
	} else {
		e.log.Warn("topic json unmarshal failed", zap.Error(err), zap.String("raw", content))
	}

	return
}
