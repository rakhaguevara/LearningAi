package infrastructure

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
)

type JWTService struct {
	accessSecret string
}

func NewJWTService(accessSecret string) *JWTService {
	return &JWTService{
		accessSecret: accessSecret,
	}
}

func (s *JWTService) GenerateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	// 1. Generate Access Token (JWT) - 15 minute expiry
	accessClaims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// 2. Generate Refresh Token (Opaque secure random string, not a JWT, since we hash it in DB)
	// Alternatively, we embed the refresh token ID into a signed JWT as requested by the use case logic.
	// Since my RefreshTokenUseCase FallbackHashScan assumes rawToken is a JWT where Subject = TokenID:
	tokenID := uuid.New().String()
	refreshClaims := jwt.MapClaims{
		"sub": tokenID,
		"usr": user.ID.String(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// We sign the refresh token payload using the same or different secret. Using accessSecret for simplicity.
	// We can use a unique REFRESH_SECRET in a real prod env.
	refreshStr, err := refreshToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr, // raw JWT string
	}, nil
}

func (s *JWTService) ValidateAccessToken(ctx context.Context, token string) (uuid.UUID, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("subject missing in token")
	}

	return uuid.Parse(sub)
}

func (s *JWTService) HashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *JWTService) CompareRefreshToken(hash, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
}

// GenerateOpaqueToken is a generic helper for opaque secrets if not using JWTs
func (s *JWTService) GenerateOpaqueToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
