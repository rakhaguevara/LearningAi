package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/adaptive-ai-learn/backend/internal/models"
	jwtpkg "github.com/adaptive-ai-learn/backend/pkg/jwt"
)

type Service struct {
	repo *Repository
	jwt  *jwtpkg.Service
	log  *zap.Logger
}

type GoogleTokenPayload struct {
	Token string `json:"token" binding:"required"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func NewService(repo *Repository, jwt *jwtpkg.Service, log *zap.Logger) *Service {
	return &Service{repo: repo, jwt: jwt, log: log}
}

func (s *Service) AuthenticateWithGoogle(idToken string) (*AuthResponse, error) {
	userInfo, err := s.verifyGoogleToken(idToken)
	if err != nil {
		return nil, apperr.NewUnauthorized("invalid google token")
	}

	if !userInfo.VerifiedEmail {
		return nil, apperr.NewUnauthorized("google email not verified")
	}

	user, err := s.repo.FindByGoogleID(userInfo.ID)
	if err == sql.ErrNoRows {
		user, err = s.repo.Create(userInfo.Email, userInfo.Name, userInfo.Picture, userInfo.ID)
		if err != nil {
			s.log.Error("failed to create user", zap.Error(err))
			return nil, apperr.NewInternal("failed to create user account")
		}

		if profileErr := s.repo.CreateLearningProfile(user.ID); profileErr != nil {
			s.log.Error("failed to create learning profile", zap.Error(profileErr))
		}

		s.log.Info("new user registered", zap.String("email", user.Email))
	} else if err != nil {
		s.log.Error("database error during auth", zap.Error(err))
		return nil, apperr.NewInternal("authentication failed")
	}

	if err := s.repo.UpdateLastLogin(user.ID); err != nil {
		s.log.Warn("failed to update last login", zap.Error(err))
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperr.NewInternal("failed to generate access token")
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperr.NewInternal("failed to generate refresh token")
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) verifyGoogleToken(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", accessToken))
	if err != nil {
		return nil, fmt.Errorf("requesting google userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google returned status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decoding google response: %w", err)
	}

	return &userInfo, nil
}

// RegisterRequest represents email/password registration payload
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents email/password login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register creates a new user with email and password
func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, apperr.NewConflict("email already registered")
	}
	if err != sql.ErrNoRows {
		s.log.Error("database error during registration", zap.Error(err))
		return nil, apperr.NewInternal("registration failed")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, apperr.NewInternal("registration failed")
	}

	// Create user
	user, err := s.repo.CreateWithPassword(req.Email, req.Name, "", string(passwordHash))
	if err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return nil, apperr.NewInternal("failed to create user account")
	}

	// Create learning profile
	if profileErr := s.repo.CreateLearningProfile(user.ID); profileErr != nil {
		s.log.Error("failed to create learning profile", zap.Error(profileErr))
	}

	s.log.Info("new user registered", zap.String("email", user.Email))

	// Generate tokens
	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperr.NewInternal("failed to generate access token")
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperr.NewInternal("failed to generate refresh token")
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates user with email and password
func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperr.NewUnauthorized("invalid email or password")
		}
		s.log.Error("database error during login", zap.Error(err))
		return nil, apperr.NewInternal("login failed")
	}

	// Verify password
	if user.PasswordHash == "" {
		return nil, apperr.NewUnauthorized("please use Google login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperr.NewUnauthorized("invalid email or password")
	}

	// Update last login
	if err := s.repo.UpdateLastLogin(user.ID); err != nil {
		s.log.Warn("failed to update last login", zap.Error(err))
	}

	// Generate tokens
	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperr.NewInternal("failed to generate access token")
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperr.NewInternal("failed to generate refresh token")
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
