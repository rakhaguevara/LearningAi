package onboarding

import (
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo *Repository
	log  *zap.Logger
}

func NewService(repo *Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

// GetStatus returns whether the user's onboarding is complete.
func (s *Service) GetStatus(userID uuid.UUID) (bool, error) {
	completed, err := s.repo.GetStatus(userID)
	if err != nil {
		s.log.Error("failed to get onboarding status", zap.String("user_id", userID.String()), zap.Error(err))
		return false, fmt.Errorf("failed to get onboarding status")
	}
	return completed, nil
}

// Submit persists the onboarding answers and marks the profile as completed.
func (s *Service) Submit(userID uuid.UUID, req OnboardingSubmitRequest) error {
	prefs := PreferencesPayload{
		LearningStyle:       req.LearningStyle,
		LongContentBehavior: req.LongContentBehavior,
		ExplanationFormat:   req.ExplanationFormat,
		InterestThemes:      req.InterestThemes,
		AnalogyTheme:        req.AnalogyTheme,
		DepthPreference:     req.DepthPreference,
		AIRetryPreference:   req.AIRetryPreference,
		StudyFocus:          req.StudyFocus,
		FileUploadHabit:     req.FileUploadHabit,
		LearningGoal:        req.LearningGoal,
	}

	if err := s.repo.UpsertLearningProfile(userID, prefs); err != nil {
		s.log.Error("failed to upsert learning profile", zap.Error(err))
		return fmt.Errorf("failed to save learning profile")
	}

	// Record implicit behavioral signals derived from onboarding answers.
	signals := map[string]interface{}{
		"preferred_depth":         req.DepthPreference,
		"long_content_tolerance":  req.LongContentBehavior,
		"file_upload_willingness": req.FileUploadHabit,
		"ai_interaction_style":    req.AIRetryPreference,
		"primary_interests":       req.InterestThemes,
	}
	if err := s.repo.InsertBehaviorSignals(userID, signals); err != nil {
		// Non-fatal: log and continue
		s.log.Warn("failed to insert behavior signals", zap.Error(err))
	}

	if err := s.repo.MarkProfileCompleted(userID); err != nil {
		s.log.Error("failed to mark profile completed", zap.Error(err))
		return fmt.Errorf("failed to mark profile as completed")
	}

	s.log.Info("onboarding completed", zap.String("user_id", userID.String()))
	return nil
}

// UpdateLearning applies partial updates to an existing learning profile.
func (s *Service) UpdateLearning(userID uuid.UUID, req UpdateLearningRequest) error {
	// Retrieve current profile to merge on top of it.
	current, err := s.repo.GetLearningProfile(userID)
	if err != nil {
		s.log.Error("failed to get current learning profile", zap.Error(err))
		return fmt.Errorf("failed to retrieve learning profile")
	}
	if current == nil {
		current = &PreferencesPayload{}
	}

	if req.LearningStyle != nil {
		current.LearningStyle = *req.LearningStyle
	}
	if req.DepthPreference != nil {
		current.DepthPreference = *req.DepthPreference
	}
	if req.InterestThemes != nil {
		current.InterestThemes = req.InterestThemes
	}
	if req.ExplanationFormat != nil {
		current.ExplanationFormat = *req.ExplanationFormat
	}
	if req.LearningGoal != nil {
		current.LearningGoal = *req.LearningGoal
	}

	if err := s.repo.UpsertLearningProfile(userID, *current); err != nil {
		s.log.Error("failed to update learning profile", zap.Error(err))
		return fmt.Errorf("failed to update learning profile")
	}

	return nil
}
