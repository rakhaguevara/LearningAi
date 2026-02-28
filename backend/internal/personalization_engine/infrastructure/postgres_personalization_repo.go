package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/adaptive-ai-learn/backend/internal/personalization_engine/domain"
)

type PostgresPersonalizationRepository struct {
	db *sql.DB
}

func NewPostgresPersonalizationRepository(db *sql.DB) *PostgresPersonalizationRepository {
	return &PostgresPersonalizationRepository{db: db}
}

func (r *PostgresPersonalizationRepository) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.UserLearningProfile, error) {
	query := `
		SELECT user_id, learning_style_score, interest_score, adaptability_index, last_updated
		FROM user_learning_profiles
		WHERE user_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)

	var profile domain.UserLearningProfile
	var styleBytes, interestBytes []byte

	err := row.Scan(
		&profile.UserID,
		&styleBytes,
		&interestBytes,
		&profile.AdaptabilityIndex,
		&profile.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile not found for user %s", userID)
		}
		return nil, err
	}

	if err := json.Unmarshal(styleBytes, &profile.LearningStyleScore); err != nil {
		profile.LearningStyleScore = make(map[string]float64) // fallback
	}
	if err := json.Unmarshal(interestBytes, &profile.InterestScore); err != nil {
		profile.InterestScore = make(map[string]float64) // fallback
	}

	return &profile, nil
}

func (r *PostgresPersonalizationRepository) SaveUserProfile(ctx context.Context, profile *domain.UserLearningProfile) error {
	styleBytes, _ := json.Marshal(profile.LearningStyleScore)
	interestBytes, _ := json.Marshal(profile.InterestScore)

	query := `
		INSERT INTO user_learning_profiles (user_id, learning_style_score, interest_score, adaptability_index, last_updated)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			learning_style_score = EXCLUDED.learning_style_score,
			interest_score = EXCLUDED.interest_score,
			adaptability_index = EXCLUDED.adaptability_index,
			last_updated = EXCLUDED.last_updated
	`
	_, err := r.db.ExecContext(ctx, query,
		profile.UserID,
		styleBytes,
		interestBytes,
		profile.AdaptabilityIndex,
		profile.LastUpdated,
	)

	return err
}

func (r *PostgresPersonalizationRepository) SaveLearningSignal(ctx context.Context, signal *domain.LearningSignal) error {
	query := `
		INSERT INTO learning_signals (id, user_id, session_id, time_spent, explanation_type, theme_used, engagement_score, feedback_score, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		signal.ID,
		signal.UserID,
		signal.SessionID,
		signal.TimeSpent,
		signal.ExplanationType,
		signal.ThemeUsed,
		signal.EngagementScore,
		signal.FeedbackScore,
		signal.CreatedAt,
	)
	return err
}

func (r *PostgresPersonalizationRepository) GetRecentSignals(ctx context.Context, userID uuid.UUID, limit int) ([]domain.LearningSignal, error) {
	query := `
		SELECT id, user_id, session_id, time_spent, explanation_type, theme_used, engagement_score, feedback_score, created_at
		FROM learning_signals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []domain.LearningSignal
	for rows.Next() {
		var s domain.LearningSignal
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.SessionID,
			&s.TimeSpent,
			&s.ExplanationType,
			&s.ThemeUsed,
			&s.EngagementScore,
			&s.FeedbackScore,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		signals = append(signals, s)
	}
	return signals, nil
}
