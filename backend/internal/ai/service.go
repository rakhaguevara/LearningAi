package ai

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	provider AIProvider
	log      *zap.Logger
}

func NewService(provider AIProvider, log *zap.Logger) *Service {
	return &Service{provider: provider, log: log}
}

func (s *Service) Explain(ctx context.Context, req ExplainRequest) (*ExplainResponse, error) {
	if len(req.Interests) == 0 {
		req.Interests = []string{"general"}
	}
	if req.Style == "" {
		req.Style = "adaptive"
	}
	if req.Difficulty == "" {
		req.Difficulty = "intermediate"
	}

	return s.provider.ExplainConcept(ctx, req)
}

func (s *Service) GenerateIllustration(ctx context.Context, req IllustrationRequest) (*IllustrationResponse, error) {
	if req.Style == "" {
		req.Style = "educational"
	}

	return s.provider.GenerateIllustration(ctx, req)
}

func (s *Service) AdaptStyle(ctx context.Context, req StyleRequest) (*StyleResponse, error) {
	return s.provider.AdaptTeachingStyle(ctx, req)
}
