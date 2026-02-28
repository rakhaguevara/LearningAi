package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IP       string `json:"-"`
	Device   string `json:"-"`
	Location string `json:"-"`
}

type LoginUseCase struct {
	repo         domain.AuthRepository
	tokenService domain.TokenService
}

func NewLoginUseCase(repo domain.AuthRepository, tokenService domain.TokenService) *LoginUseCase {
	return &LoginUseCase{
		repo:         repo,
		tokenService: tokenService,
	}
}

// Execute performs the login, creates the session tokens, and logs the history.
func (u *LoginUseCase) Execute(ctx context.Context, req LoginReq) (*domain.TokenPair, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("missing email or password")
	}

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.PasswordHash == nil {
		return nil, errors.New("invalid provider for this account")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate Tokens
	tokenPair, err := u.tokenService.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	// Hash and store the refresh token
	hash, err := u.tokenService.HashRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, err
	}

	rt := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}
	if err := u.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	// Log login history
	history := &domain.LoginHistory{
		ID:        uuid.New(),
		UserID:    user.ID,
		IP:        req.IP,
		Device:    req.Device,
		Location:  req.Location,
		CreatedAt: time.Now(),
	}
	_ = u.repo.CreateLoginHistory(ctx, history)

	return tokenPair, nil
}
