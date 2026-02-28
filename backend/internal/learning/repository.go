package learning

import (
	"database/sql"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/adaptive-ai-learn/backend/internal/models"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) CreateSession(userID uuid.UUID, topic, subject, style string) (*models.LearningSession, error) {
	var s models.LearningSession
	err := r.db.QueryRow(`
		INSERT INTO learning_sessions (user_id, topic, subject, style, status)
		VALUES ($1, $2, $3, $4, 'active')
		RETURNING id, user_id, topic, subject, style, status, started_at, ended_at, duration_sec, created_at
	`, userID, topic, subject, style).Scan(
		&s.ID, &s.UserID, &s.Topic, &s.Subject, &s.Style,
		&s.Status, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetSession(sessionID uuid.UUID) (*models.LearningSession, error) {
	var s models.LearningSession
	err := r.db.QueryRow(`
		SELECT id, user_id, topic, subject, style, status, started_at, ended_at, duration_sec, created_at
		FROM learning_sessions WHERE id = $1
	`, sessionID).Scan(
		&s.ID, &s.UserID, &s.Topic, &s.Subject, &s.Style,
		&s.Status, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) GetUserSessions(userID uuid.UUID, limit int) ([]models.LearningSession, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, topic, subject, style, status, started_at, ended_at, duration_sec, created_at
		FROM learning_sessions WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.LearningSession
	for rows.Next() {
		var s models.LearningSession
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.Topic, &s.Subject, &s.Style,
			&s.Status, &s.StartedAt, &s.EndedAt, &s.DurationSec, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}

func (r *Repository) EndSession(sessionID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE learning_sessions
		SET status = 'completed',
		    ended_at = NOW(),
		    duration_sec = EXTRACT(EPOCH FROM (NOW() - started_at))::INT
		WHERE id = $1
	`, sessionID)
	return err
}

func (r *Repository) SaveInteraction(interaction *models.AIInteractionHistory) error {
	_, err := r.db.Exec(`
		INSERT INTO ai_interaction_history (session_id, user_id, prompt, response, interaction_type, tokens_used, latency_ms)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, interaction.SessionID, interaction.UserID, interaction.Prompt, interaction.Response,
		interaction.InteractionType, interaction.TokensUsed, interaction.LatencyMs)
	return err
}
