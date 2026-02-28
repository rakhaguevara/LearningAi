package user

import (
	"database/sql"

	"github.com/google/uuid"
	"go.uber.org/zap"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/adaptive-ai-learn/backend/internal/models"
)

type Service struct {
	repo *Repository
	log  *zap.Logger
}

type ProfileResponse struct {
	User            *models.User            `json:"user"`
	LearningProfile *models.LearningProfile `json:"learning_profile,omitempty"`
	Interests       []models.InterestTag    `json:"interests"`
}

func NewService(repo *Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) GetProfile(userID uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err == sql.ErrNoRows {
		return nil, apperr.NewNotFound("user not found")
	}
	if err != nil {
		s.log.Error("failed to fetch user", zap.Error(err))
		return nil, apperr.NewInternal("failed to fetch user profile")
	}

	profile, err := s.repo.GetLearningProfile(userID)
	if err != nil && err != sql.ErrNoRows {
		s.log.Error("failed to fetch learning profile", zap.Error(err))
	}

	interests, err := s.repo.GetInterestTags(userID)
	if err != nil {
		s.log.Error("failed to fetch interests", zap.Error(err))
		interests = []models.InterestTag{}
	}

	return &ProfileResponse{
		User:            user,
		LearningProfile: profile,
		Interests:       interests,
	}, nil
}
