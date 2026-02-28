package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserLearningProfile represents a user's accumulated learning preferences.
type UserLearningProfile struct {
	UserID             uuid.UUID          `json:"user_id"`
	LearningStyleScore map[string]float64 `json:"learning_style_score"`
	InterestScore      map[string]float64 `json:"interest_score"`
	AdaptabilityIndex  float64            `json:"adaptability_index"`
	LastUpdated        time.Time          `json:"last_updated"`
}

// LearningSignal represents a single atomic interaction event.
type LearningSignal struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	SessionID       uuid.UUID `json:"session_id"`
	TimeSpent       int       `json:"time_spent"` // in seconds
	ExplanationType string    `json:"explanation_type"`
	ThemeUsed       string    `json:"theme_used"`
	EngagementScore float64   `json:"engagement_score"` // 0.0 to 1.0 (e.g., scroll depth, completion)
	FeedbackScore   float64   `json:"feedback_score"`   // 0.0 to 1.0 (e.g., explicit thumbs up/down, or quiz pass rate)
	CreatedAt       time.Time `json:"created_at"`
}
