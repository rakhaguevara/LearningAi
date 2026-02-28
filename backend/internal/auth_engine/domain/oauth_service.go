package domain

import (
	"context"
)

// GoogleUserData represents the generic user info fetched from Google.
type GoogleUserData struct {
	GoogleID  string
	Email     string
	Name      string
	AvatarURL string
}

// OAuthService defines operations for external OAuth providers like Google.
type OAuthService interface {
	// GetLoginURL returns the redirect URL for the OAuth consent screen.
	GetLoginURL(state string) string

	// ExchangeCodeForUser takes the auth code and retrieves the user profile from Google.
	ExchangeCodeForUser(ctx context.Context, code string) (*GoogleUserData, error)
}
