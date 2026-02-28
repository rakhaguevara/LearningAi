package learning

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperr "github.com/adaptive-ai-learn/backend/internal/common/errors"
	"github.com/adaptive-ai-learn/backend/internal/models"
)

type Service struct {
	repo *Repository
	log  *zap.Logger
}

type StartSessionRequest struct {
	Topic   string `json:"topic" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Style   string `json:"style"`
}

type StartSessionResponse struct {
	Session *models.LearningSession `json:"session"`
}

func NewService(repo *Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) StartSession(userID uuid.UUID, req StartSessionRequest) (*StartSessionResponse, error) {
	if req.Style == "" {
		req.Style = "adaptive"
	}

	session, err := s.repo.CreateSession(userID, req.Topic, req.Subject, req.Style)
	if err != nil {
		s.log.Error("failed to create learning session", zap.Error(err))
		return nil, apperr.NewInternal("failed to start learning session")
	}

	s.log.Info("learning session started",
		zap.String("session_id", session.ID.String()),
		zap.String("topic", req.Topic),
	)

	return &StartSessionResponse{Session: session}, nil
}

func (s *Service) GetUserSessions(userID uuid.UUID) ([]models.LearningSession, error) {
	sessions, err := s.repo.GetUserSessions(userID, 50)
	if err != nil {
		s.log.Error("failed to fetch sessions", zap.Error(err))
		return nil, apperr.NewInternal("failed to fetch learning sessions")
	}
	if sessions == nil {
		sessions = []models.LearningSession{}
	}
	return sessions, nil
}
