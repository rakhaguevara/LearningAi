package usecase_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/adaptive-ai-learn/backend/internal/personalization_engine/domain"
	"github.com/adaptive-ai-learn/backend/internal/personalization_engine/usecase"
)

// MockRepo for testing demonstration
type MockRepo struct {
	profile *domain.UserLearningProfile
}

func (m *MockRepo) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.UserLearningProfile, error) {
	if m.profile != nil {
		return m.profile, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockRepo) SaveUserProfile(ctx context.Context, profile *domain.UserLearningProfile) error {
	return nil
}
func (m *MockRepo) SaveLearningSignal(ctx context.Context, signal *domain.LearningSignal) error {
	return nil
}
func (m *MockRepo) GetRecentSignals(ctx context.Context, userID uuid.UUID, limit int) ([]domain.LearningSignal, error) {
	return nil, nil
}

func ExamplePersonalizeExplanationUseCase_Execute() {
	// Setup mock profile where user loves anime and visual storytelling
	userID := uuid.New()
	mockProfile := &domain.UserLearningProfile{
		UserID: userID,
		LearningStyleScore: map[string]float64{
			"visual":    0.8,
			"concise":   0.1,
			"narrative": 0.1,
		},
		InterestScore: map[string]float64{
			"anime":      0.9,
			"sports":     0.05,
			"technology": 0.05,
		},
	}

	repo := &MockRepo{profile: mockProfile}
	styleClassifier := domain.NewRuleBasedLearningStyleClassifier()
	interestClassifier := domain.NewRuleBasedInterestClassifier()

	engine := domain.NewPersonalizationEngine(styleClassifier, interestClassifier, repo)
	uc := usecase.NewPersonalizeExplanationUseCase(engine)

	basePrompt := "Explain the concept of quantum entanglement."

	adaptedPrompt, err := uc.Execute(context.Background(), userID.String(), basePrompt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("=== Base Prompt ===")
	fmt.Println(basePrompt)
	fmt.Println("\n=== Adapted Prompt ===")
	fmt.Println(adaptedPrompt)

	// Output:
	// === Base Prompt ===
	// Explain the concept of quantum entanglement.
	//
	// === Adapted Prompt ===
	// You are an adaptive AI tutor. The student prefers a visual approach and is highly interested in anime. Integrate anime analogies and vivid visual themes into your explanation.
	//
	// Original Request:
	// Explain the concept of quantum entanglement.
}
