package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user.
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash *string   `json:"-"`        // nullable, hidden from JSON
	Provider     string    `json:"provider"` // enum: "local", "google"
	ProviderID   *string   `json:"provider_id"`
	Name         string    `json:"name"`
	AvatarURL    *string   `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RefreshToken represents a long-lived session token tied to a user.
type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"` // bcrypt hash of the plain token
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginHistory represents a record of a single login event.
type LoginHistory struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	IP        string    `json:"ip"`
	Device    string    `json:"device"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthProvider constants
const (
	AuthProviderLocal  = "local"
	AuthProviderGoogle = "google"
)
