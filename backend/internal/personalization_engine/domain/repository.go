package domain

import (
	"context"

	"github.com/google/uuid"
)

// PersonalizationRepository defines the contract for persisting personalization data.
type PersonalizationRepository interface {
	// GetUserProfile retrieves the learning profile for a given user.
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserLearningProfile, error)

	// SaveUserProfile stores or updates a user's learning profile.
	SaveUserProfile(ctx context.Context, profile *UserLearningProfile) error

	// SaveLearningSignal records a new learning interaction signal.
	SaveLearningSignal(ctx context.Context, signal *LearningSignal) error

	// GetRecentSignals retrieves the N most recent signals for a user to aid in classification.
	GetRecentSignals(ctx context.Context, userID uuid.UUID, limit int) ([]LearningSignal, error)
}
