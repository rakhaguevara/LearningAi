package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/adaptive-ai-learn/backend/internal/config"
)

func Connect(cfg config.DBConfig, log *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connectivity with retries
	var pingErr error
	for i := 0; i < 5; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		log.Warn("database ping failed, retrying...", zap.Int("attempt", i+1), zap.Error(pingErr))
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		return nil, fmt.Errorf("pinging database after retries: %w", pingErr)
	}

	log.Info("database connection established",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.Name),
	)

	return db, nil
}

func RunMigrations(db *sql.DB, log *zap.Logger) error {
	migrations := []string{
		createUsersTable,
		createLearningProfilesTable,
		createInterestTagsTable,
		createLearningSessionsTable,
		createAIInteractionHistoryTable,
		createBehaviorSignalsTable,
		createLearningStyleProfilesTable,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("running migration %d: %w", i, err)
		}
	}

	log.Info("database migrations completed")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    avatar_url    TEXT DEFAULT '',
    google_id     VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    role          VARCHAR(50) DEFAULT 'learner',
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
`

const createLearningProfilesTable = `
CREATE TABLE IF NOT EXISTS learning_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    preferred_style     VARCHAR(100) DEFAULT 'adaptive',
    difficulty_level    VARCHAR(50) DEFAULT 'intermediate',
    learning_goals      TEXT[] DEFAULT '{}',
    weekly_target_hours FLOAT DEFAULT 5.0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);
`

const createInterestTagsTable = `
CREATE TABLE IF NOT EXISTS interest_tags (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag        VARCHAR(100) NOT NULL,
    category   VARCHAR(100) NOT NULL,
    weight     FLOAT DEFAULT 1.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, tag)
);
CREATE INDEX IF NOT EXISTS idx_interest_tags_user ON interest_tags(user_id);
`

const createLearningSessionsTable = `
CREATE TABLE IF NOT EXISTS learning_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic        VARCHAR(500) NOT NULL,
    subject      VARCHAR(255) NOT NULL,
    style        VARCHAR(100) DEFAULT 'adaptive',
    status       VARCHAR(50) DEFAULT 'active',
    started_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at     TIMESTAMPTZ,
    duration_sec INT DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_learning_sessions_user ON learning_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_sessions_status ON learning_sessions(status);
`

const createAIInteractionHistoryTable = `
CREATE TABLE IF NOT EXISTS ai_interaction_history (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id       UUID NOT NULL REFERENCES learning_sessions(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prompt           TEXT NOT NULL,
    response         TEXT NOT NULL,
    interaction_type VARCHAR(100) NOT NULL,
    tokens_used      INT DEFAULT 0,
    latency_ms       INT DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ai_history_session ON ai_interaction_history(session_id);
CREATE INDEX IF NOT EXISTS idx_ai_history_user ON ai_interaction_history(user_id);
`

const createBehaviorSignalsTable = `
CREATE TABLE IF NOT EXISTS behavior_signals (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id  UUID REFERENCES learning_sessions(id) ON DELETE SET NULL,
    signal_type VARCHAR(100) NOT NULL,
    value       FLOAT DEFAULT 1.0,
    context     JSONB DEFAULT '{}',
    topic       VARCHAR(500),
    subject     VARCHAR(255),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_behavior_signals_user ON behavior_signals(user_id);
CREATE INDEX IF NOT EXISTS idx_behavior_signals_type ON behavior_signals(signal_type);
CREATE INDEX IF NOT EXISTS idx_behavior_signals_created ON behavior_signals(created_at DESC);
`

const createLearningStyleProfilesTable = `
CREATE TABLE IF NOT EXISTS learning_style_profiles (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    primary_style      VARCHAR(50) DEFAULT 'adaptive',
    visual_score       FLOAT DEFAULT 0.25,
    auditory_score     FLOAT DEFAULT 0.25,
    reading_score      FLOAT DEFAULT 0.25,
    kinesthetic_score  FLOAT DEFAULT 0.25,
    confidence         FLOAT DEFAULT 0,
    sample_size        INT DEFAULT 0,
    last_calculated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);
CREATE INDEX IF NOT EXISTS idx_learning_style_user ON learning_style_profiles(user_id);
`
