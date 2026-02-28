package usecase

import (
	"context"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
	"github.com/google/uuid"
)

type LogoutUseCase struct {
	repo domain.AuthRepository
}

func NewLogoutUseCase(repo domain.AuthRepository) *LogoutUseCase {
	return &LogoutUseCase{repo: repo}
}

// Execute revokes all tokens for a user, forcing them to login again
func (u *LogoutUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	// A more granular logout might only revoke the current token.
	// But deleting all tokens for the user is safer and simpler for "logout everywhere".
	return u.repo.DeleteRefreshTokensByUserID(ctx, userID)
}
