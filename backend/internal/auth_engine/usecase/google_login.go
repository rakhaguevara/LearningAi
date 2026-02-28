package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
)

type GoogleOAuthUseCase struct {
	repo         domain.AuthRepository
	tokenService domain.TokenService
	oauthService domain.OAuthService
}

func NewGoogleOAuthUseCase(repo domain.AuthRepository, tokenService domain.TokenService, oauthService domain.OAuthService) *GoogleOAuthUseCase {
	return &GoogleOAuthUseCase{
		repo:         repo,
		tokenService: tokenService,
		oauthService: oauthService,
	}
}

// GetLoginURL returns the formatted redirect URL to Google
func (u *GoogleOAuthUseCase) GetLoginURL(state string) string {
	return u.oauthService.GetLoginURL(state)
}

// HandleCallback exchanges the code, creates/finds the user, and generates tokens
func (u *GoogleOAuthUseCase) HandleCallback(ctx context.Context, code, ip, device string) (*domain.TokenPair, error) {
	if code == "" {
		return nil, errors.New("auth code is required")
	}

	// Exchange code for Google user profile
	googleUser, err := u.oauthService.ExchangeCodeForUser(ctx, code)
	if err != nil {
		return nil, err
	}

	// Find or Create User
	user, err := u.repo.GetUserByEmail(ctx, googleUser.Email)
	if err != nil {
		// Create new user since they don't exist
		user = &domain.User{
			ID:         uuid.New(),
			Email:      googleUser.Email,
			Provider:   domain.AuthProviderGoogle,
			ProviderID: &googleUser.GoogleID,
			Name:       googleUser.Name,
			AvatarURL:  &googleUser.AvatarURL,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := u.repo.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	} else if user.Provider != domain.AuthProviderGoogle {
		// For simplicity, we could merge accounts or return an error. Returning error for security.
		return nil, errors.New("email already linked to a native account, please login using password")
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

	// Log Login History
	history := &domain.LoginHistory{
		ID:        uuid.New(),
		UserID:    user.ID,
		IP:        ip,
		Device:    device,
		Location:  "Google OAuth",
		CreatedAt: time.Now(),
	}
	_ = u.repo.CreateLoginHistory(ctx, history)

	return tokenPair, nil
}
