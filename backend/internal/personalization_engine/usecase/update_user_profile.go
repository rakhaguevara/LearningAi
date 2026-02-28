package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/adaptive-ai-learn/backend/internal/personalization_engine/domain"
)

type UpdateUserProfileUseCase struct {
	repo               domain.PersonalizationRepository
	styleClassifier    domain.LearningStyleClassifier
	interestClassifier domain.InterestClassifier
}

func NewUpdateUserProfileUseCase(
	repo domain.PersonalizationRepository,
	styleClassifier domain.LearningStyleClassifier,
	interestClassifier domain.InterestClassifier,
) *UpdateUserProfileUseCase {
	return &UpdateUserProfileUseCase{
		repo:               repo,
		styleClassifier:    styleClassifier,
		interestClassifier: interestClassifier,
	}
}

// Execute processes a new learning signal, saves it, and recalculates the user's profile.
func (uc *UpdateUserProfileUseCase) Execute(ctx context.Context, signal *domain.LearningSignal) error {
	// 1. Save the new signal
	if signal.ID == uuid.Nil {
		signal.ID = uuid.New()
	}
	if signal.CreatedAt.IsZero() {
		signal.CreatedAt = time.Now()
	}

	if err := uc.repo.SaveLearningSignal(ctx, signal); err != nil {
		return fmt.Errorf("failed to save learning signal: %w", err)
	}

	// 2. Fetch the user's current profile. If not exists, create a new one.
	profile, err := uc.repo.GetUserProfile(ctx, signal.UserID)
	if err != nil { // assume not found error handled in repo by returning an empty profile
		profile = &domain.UserLearningProfile{
			UserID:             signal.UserID,
			LearningStyleScore: make(map[string]float64),
			InterestScore:      make(map[string]float64),
		}
	}

	// 3. Fetch recent signals to recalculate dominant traits (e.g. trailing 50 signals)
	recentSignals, err := uc.repo.GetRecentSignals(ctx, signal.UserID, 50)
	if err != nil {
		return fmt.Errorf("failed to fetch recent signals: %w", err)
	}

	// 4. Time Decay existing profile scores before adding new classifications
	daysSinceUpdate := int(time.Since(profile.LastUpdated).Hours() / 24)
	if daysSinceUpdate > 0 {
		for k, v := range profile.LearningStyleScore {
			profile.LearningStyleScore[k] = domain.ApplyTimeDecay(v, daysSinceUpdate)
		}
		for k, v := range profile.InterestScore {
			profile.InterestScore[k] = domain.ApplyTimeDecay(v, daysSinceUpdate)
		}
	}

	// 5. Run classifiers on recent signals
	newStyleScores := uc.styleClassifier.Classify(recentSignals)
	newInterestScores := uc.interestClassifier.Classify(recentSignals)

	// 6. Merge new classifications into profile (simple moving average/additive approach)
	for k, v := range newStyleScores {
		profile.LearningStyleScore[k] += v
	}
	for k, v := range newInterestScores {
		profile.InterestScore[k] += v
	}

	// Re-normalize profile scores after merge to keep them bounded
	profile.LearningStyleScore = normalizeScores(profile.LearningStyleScore)
	profile.InterestScore = normalizeScores(profile.InterestScore)

	profile.LastUpdated = time.Now()

	// 7. Save updated profile
	if err := uc.repo.SaveUserProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

func normalizeScores(scores map[string]float64) map[string]float64 {
	var total float64
	for _, v := range scores {
		total += v
	}

	if total == 0 {
		return scores
	}

	normalized := make(map[string]float64)
	for k, v := range scores {
		normalized[k] = v / total
	}
	return normalized
}
