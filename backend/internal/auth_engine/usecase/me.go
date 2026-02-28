package usecase

import (
	"context"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url"`
	Provider  string    `json:"provider"`
}

type MeUseCase struct {
	repo domain.AuthRepository
}

func NewMeUseCase(repo domain.AuthRepository) *MeUseCase {
	return &MeUseCase{repo: repo}
}

func (u *MeUseCase) Execute(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Provider:  user.Provider,
	}, nil
}
