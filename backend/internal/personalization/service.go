package personalization

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service provides the main API for the personalization engine.
type Service struct {
	repo               *Repository
	styleClassifier    *LearningStyleClassifier
	interestClassifier *InterestClassifier
	scoringEngine      *ScoringEngine
	toneGenerator      *ToneGenerator
	log                *zap.Logger
}

func NewService(repo *Repository, log *zap.Logger) *Service {
	styleClassifier := NewLearningStyleClassifier(repo, log)
	interestClassifier := NewInterestClassifier(repo, log)
	scoringEngine := NewScoringEngine(repo, styleClassifier, interestClassifier, log)
	toneGenerator := NewToneGenerator()

	return &Service{
		repo:               repo,
		styleClassifier:    styleClassifier,
		interestClassifier: interestClassifier,
		scoringEngine:      scoringEngine,
		toneGenerator:      toneGenerator,
		log:                log,
	}
}

// GetPersonalizationProfile returns the complete personalization profile for a user.
func (s *Service) GetPersonalizationProfile(ctx context.Context, userID uuid.UUID) (*PersonalizationProfile, error) {
	return s.scoringEngine.BuildProfile(userID)
}

// GetAdaptivePrompt generates an AI system prompt tailored to the user.
func (s *Service) GetAdaptivePrompt(ctx context.Context, userID uuid.UUID) (string, error) {
	profile, err := s.scoringEngine.BuildProfile(userID)
	if err != nil {
		return "", err
	}
	return profile.AdaptivePrompt, nil
}

// GetToneConfig returns numeric tone parameters for the user.
func (s *Service) GetToneConfig(ctx context.Context, userID uuid.UUID) (*ToneConfig, error) {
	profile, err := s.scoringEngine.BuildProfile(userID)
	if err != nil {
		return nil, err
	}
	config := s.toneGenerator.GenerateToneConfig(profile)
	return &config, nil
}

// RecordBehaviorSignal records a user behavior signal for personalization.
func (s *Service) RecordBehaviorSignal(ctx context.Context, signal *BehaviorSignal) error {
	return s.repo.RecordSignal(signal)
}

// RecordExplanationRequest records that a user requested an explanation.
func (s *Service) RecordExplanationRequest(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID, topic, subject string) error {
	signal := &BehaviorSignal{
		UserID:     userID,
		SessionID:  sessionID,
		SignalType: SignalExplanationRequest,
		Value:      1.0,
		Topic:      topic,
		Subject:    subject,
		Context:    map[string]interface{}{},
	}
	return s.repo.RecordSignal(signal)
}

// RecordIllustrationRequest records that a user requested an illustration.
func (s *Service) RecordIllustrationRequest(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID, topic string) error {
	signal := &BehaviorSignal{
		UserID:     userID,
		SessionID:  sessionID,
		SignalType: SignalIllustrationRequest,
		Value:      1.0,
		Topic:      topic,
		Context:    map[string]interface{}{},
	}
	return s.repo.RecordSignal(signal)
}

// RecordFollowUpQuestion records that a user asked a follow-up question.
func (s *Service) RecordFollowUpQuestion(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID, topic string) error {
	signal := &BehaviorSignal{
		UserID:     userID,
		SessionID:  sessionID,
		SignalType: SignalFollowUpQuestion,
		Value:      1.0,
		Topic:      topic,
		Context:    map[string]interface{}{},
	}
	return s.repo.RecordSignal(signal)
}

// RecordSessionDuration records the duration of a learning session.
func (s *Service) RecordSessionDuration(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID, durationSec int) error {
	// Normalize duration to 0-1 scale (cap at 30 min = 1.0)
	normalizedValue := float64(durationSec) / 1800
	if normalizedValue > 1.0 {
		normalizedValue = 1.0
	}

	signal := &BehaviorSignal{
		UserID:     userID,
		SessionID:  sessionID,
		SignalType: SignalSessionDuration,
		Value:      normalizedValue,
		Context:    map[string]interface{}{"duration_sec": durationSec},
	}
	return s.repo.RecordSignal(signal)
}

// RecordDifficultyFeedback records user feedback on content difficulty.
func (s *Service) RecordDifficultyFeedback(ctx context.Context, userID uuid.UUID, sessionID *uuid.UUID, feedback string) error {
	// Map feedback to numeric value
	var value float64
	switch feedback {
	case "too_easy":
		value = 0.2
	case "just_right":
		value = 0.5
	case "too_hard":
		value = 0.8
	default:
		value = 0.5
	}

	signal := &BehaviorSignal{
		UserID:     userID,
		SessionID:  sessionID,
		SignalType: SignalDifficultyFeedback,
		Value:      value,
		Context:    map[string]interface{}{"feedback": feedback},
	}
	return s.repo.RecordSignal(signal)
}

// RecordInterestIndication records that a user showed interest in a topic.
func (s *Service) RecordInterestIndication(ctx context.Context, userID uuid.UUID, interest string, intensity float64) error {
	return s.interestClassifier.RecordInterestSignal(userID, interest, intensity)
}

// RecordAnalogySatisfaction records how well an analogy resonated with the user.
func (s *Service) RecordAnalogySatisfaction(ctx context.Context, userID uuid.UUID, analogyDomain string, satisfaction float64) error {
	return s.interestClassifier.RecordAnalogySatisfaction(userID, analogyDomain, satisfaction)
}

// GetLearningStyleProfile returns the user's learning style classification.
func (s *Service) GetLearningStyleProfile(ctx context.Context, userID uuid.UUID) (*LearningStyleProfile, error) {
	return s.styleClassifier.Classify(userID)
}

// GetInterestProfile returns the user's interest classification.
func (s *Service) GetInterestProfile(ctx context.Context, userID uuid.UUID) (*InterestProfile, error) {
	return s.interestClassifier.Classify(userID)
}

// GetPersonalizationScore returns a 0-100 score indicating personalization depth.
func (s *Service) GetPersonalizationScore(ctx context.Context, userID uuid.UUID) (int, error) {
	profile, err := s.scoringEngine.BuildProfile(userID)
	if err != nil {
		return 0, err
	}
	return s.scoringEngine.CalculatePersonalizationScore(profile), nil
}

// AddUserInterest explicitly adds an interest for a user.
func (s *Service) AddUserInterest(ctx context.Context, userID uuid.UUID, tag, category string, weight float64) error {
	return s.repo.UpdateInterestWeight(userID, tag, category, weight)
}
