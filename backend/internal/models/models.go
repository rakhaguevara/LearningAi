package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Email            string     `json:"email" db:"email"`
	Name             string     `json:"name" db:"name"`
	AvatarURL        string     `json:"avatar_url" db:"avatar_url"`
	GoogleID         string     `json:"-" db:"google_id"`
	PasswordHash     string     `json:"-" db:"password_hash"`
	Role             string     `json:"role" db:"role"`
	ProfileCompleted bool       `json:"profile_completed" db:"profile_completed"`
	LastLoginAt      *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type LearningProfile struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	PreferredStyle    string    `json:"preferred_style" db:"preferred_style"`
	DifficultyLevel   string    `json:"difficulty_level" db:"difficulty_level"`
	LearningGoals     []string  `json:"learning_goals" db:"learning_goals"`
	WeeklyTargetHours float64   `json:"weekly_target_hours" db:"weekly_target_hours"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type InterestTag struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Tag       string    `json:"tag" db:"tag"`
	Category  string    `json:"category" db:"category"`
	Weight    float64   `json:"weight" db:"weight"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LearningSession struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Topic       string     `json:"topic" db:"topic"`
	Subject     string     `json:"subject" db:"subject"`
	Style       string     `json:"style" db:"style"`
	Status      string     `json:"status" db:"status"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	DurationSec int        `json:"duration_sec" db:"duration_sec"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

type AIInteractionHistory struct {
	ID              uuid.UUID `json:"id" db:"id"`
	SessionID       uuid.UUID `json:"session_id" db:"session_id"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	Prompt          string    `json:"prompt" db:"prompt"`
	Response        string    `json:"response" db:"response"`
	InteractionType string    `json:"interaction_type" db:"interaction_type"`
	TokensUsed      int       `json:"tokens_used" db:"tokens_used"`
	LatencyMs       int       `json:"latency_ms" db:"latency_ms"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
