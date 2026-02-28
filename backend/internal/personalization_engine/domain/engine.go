package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// PersonalizationEngine ties the domain components together to adapt system prompts.
// It relies on abstractions to allow swapping out logic (like ML classifiers) in the future.
type PersonalizationEngine struct {
	styleClassifier    LearningStyleClassifier
	interestClassifier InterestClassifier
	repo               PersonalizationRepository
}

func NewPersonalizationEngine(sc LearningStyleClassifier, ic InterestClassifier, repo PersonalizationRepository) *PersonalizationEngine {
	return &PersonalizationEngine{
		styleClassifier:    sc,
		interestClassifier: ic,
		repo:               repo,
	}
}

// AdaptPrompt fetches the user's profile, determines their dominant style & interest,
// and prepends tailored instructions to the base LLM prompt.
func (e *PersonalizationEngine) AdaptPrompt(ctx context.Context, userID string, basePrompt string) (string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	profile, err := e.repo.GetUserProfile(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("failed to get user profile: %w", err)
	}

	// Determine dominant traits
	dominantStyle := getDominantTrait(profile.LearningStyleScore)
	dominantInterest := getDominantTrait(profile.InterestScore)

	var adaptation string

	if dominantInterest != "" && dominantStyle != "" {
		adaptation = fmt.Sprintf("You are an adaptive AI tutor. The student prefers a %s approach and is highly interested in %s. Integrate %s analogies and vivid %s themes into your explanation.", dominantStyle, dominantInterest, dominantInterest, dominantStyle)
	} else if dominantInterest != "" {
		adaptation = fmt.Sprintf("You are an adaptive AI tutor. Explain the concept using analogies from the domain of %s.", dominantInterest)
	} else if dominantStyle != "" {
		adaptation = fmt.Sprintf("You are an adaptive AI tutor. Explain the student's question using a %s style.", dominantStyle)
	} else {
		adaptation = "You are a helpful and adaptive AI tutor. Provide a clear and concise explanation." // Fallback
	}

	// Prepend adaptation to base prompt
	adaptedPrompt := fmt.Sprintf("%s\n\nOriginal Request:\n%s", adaptation, basePrompt)
	return adaptedPrompt, nil
}

func getDominantTrait(scores map[string]float64) string {
	var maxKey string
	var maxVal float64 = -1.0

	for k, v := range scores {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}
	return maxKey
}
