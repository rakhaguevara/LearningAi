package domain

import (
	"context"

	"github.com/google/uuid"
)

// AuthRepository defines the database interactions for the authentication domain.
type AuthRepository interface {
	// User
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error

	// Refresh Token
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
	DeleteRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error

	// Login History
	CreateLoginHistory(ctx context.Context, history *LoginHistory) error
}
