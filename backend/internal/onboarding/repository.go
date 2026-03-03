package onboarding

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

// GetStatus returns true if the user's profile is marked completed.
func (r *Repository) GetStatus(userID uuid.UUID) (bool, error) {
	var completed bool
	err := r.db.QueryRow(
		`SELECT profile_completed FROM users WHERE id = $1`, userID,
	).Scan(&completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user not found")
		}
		return false, err
	}
	return completed, nil
}

// MarkProfileCompleted sets profile_completed = true on the users row.
func (r *Repository) MarkProfileCompleted(userID uuid.UUID) error {
	_, err := r.db.Exec(
		`UPDATE users SET profile_completed = true, updated_at = $1 WHERE id = $2`,
		time.Now(), userID,
	)
	return err
}

// UpsertLearningProfile inserts or updates the onboarding preferences JSONB.
func (r *Repository) UpsertLearningProfile(userID uuid.UUID, prefs PreferencesPayload) error {
	prefsJSON, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("marshalling preferences: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO user_learning_profiles (user_id, preferences, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET preferences = $2, updated_at = NOW()
	`, userID, prefsJSON)
	return err
}

// InsertBehaviorSignals records implicit signals gathered during onboarding.
func (r *Repository) InsertBehaviorSignals(userID uuid.UUID, signals map[string]interface{}) error {
	signalsJSON, err := json.Marshal(signals)
	if err != nil {
		return fmt.Errorf("marshalling signals: %w", err)
	}

	_, err = r.db.Exec(`
		INSERT INTO user_behavior_signals (user_id, source, signals)
		VALUES ($1, 'onboarding', $2)
	`, userID, signalsJSON)
	return err
}

// GetLearningProfile retrieves the raw preferences JSONB for a user.
func (r *Repository) GetLearningProfile(userID uuid.UUID) (*PreferencesPayload, error) {
	var raw []byte
	err := r.db.QueryRow(
		`SELECT preferences FROM user_learning_profiles WHERE user_id = $1`, userID,
	).Scan(&raw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var prefs PreferencesPayload
	if err := json.Unmarshal(raw, &prefs); err != nil {
		return nil, fmt.Errorf("unmarshalling preferences: %w", err)
	}
	return &prefs, nil
}
