package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
)

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterUseCase struct {
	repo domain.AuthRepository
}

func NewRegisterUseCase(repo domain.AuthRepository) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

func (u *RegisterUseCase) Execute(ctx context.Context, req RegisterReq) error {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return errors.New("missing required fields")
	}

	// Check if user exists
	existing, _ := u.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return errors.New("email already in use")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashStr := string(hash)

	// Create User
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: &hashStr,
		Provider:     domain.AuthProviderLocal,
		Name:         req.Name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
