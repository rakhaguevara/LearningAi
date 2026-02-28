package usecase

import (
	"context"

	"github.com/adaptive-ai-learn/backend/internal/personalization_engine/domain"
)

type PersonalizeExplanationUseCase struct {
	engine *domain.PersonalizationEngine
}

func NewPersonalizeExplanationUseCase(engine *domain.PersonalizationEngine) *PersonalizeExplanationUseCase {
	return &PersonalizeExplanationUseCase{
		engine: engine,
	}
}

// Execute adapts a raw system prompt based on the user's learned profile.
func (uc *PersonalizeExplanationUseCase) Execute(ctx context.Context, userID string, basePrompt string) (string, error) {
	return uc.engine.AdaptPrompt(ctx, userID, basePrompt)
}
