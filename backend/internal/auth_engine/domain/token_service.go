package domain

import (
	"context"

	"github.com/google/uuid"
)

// TokenPair holds the generated access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// TokenService defines operations for creating and verifying JWT and refresh tokens.
type TokenService interface {
	// GenerateTokenPair generates an access token (JWT) and a raw refresh token.
	GenerateTokenPair(ctx context.Context, user *User) (*TokenPair, error)

	// ValidateAccessToken verifies a JWT and returns the extracted user ID.
	ValidateAccessToken(ctx context.Context, token string) (uuid.UUID, error)

	// HashRefreshToken hashes a raw refresh token string for safe database storage.
	HashRefreshToken(token string) (string, error)

	// CompareRefreshToken compares a raw token with a hashed token.
	CompareRefreshToken(hash, raw string) error
}
