package auth

import (
	"database/sql"
	"time"

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

func (r *Repository) FindByGoogleID(googleID string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT id, email, name, avatar_url, google_id, role, last_login_at, created_at, updated_at
		FROM users WHERE google_id = $1
	`, googleID).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.GoogleID,
		&u.Role, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	var googleID sql.NullString
	var passwordHash sql.NullString
	err := r.db.QueryRow(`
		SELECT id, email, name, avatar_url, google_id, password_hash, role, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &googleID, &passwordHash,
		&u.Role, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if googleID.Valid {
		u.GoogleID = googleID.String
	}
	if passwordHash.Valid {
		u.PasswordHash = passwordHash.String
	}
	return &u, nil
}

func (r *Repository) Create(email, name, avatarURL, googleID string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		INSERT INTO users (email, name, avatar_url, google_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, avatar_url, google_id, role, last_login_at, created_at, updated_at
	`, email, name, avatarURL, googleID).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.GoogleID,
		&u.Role, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateWithPassword(email, name, avatarURL, passwordHash string) (*models.User, error) {
	var u models.User
	var googleID sql.NullString
	err := r.db.QueryRow(`
		INSERT INTO users (email, name, avatar_url, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, avatar_url, google_id, password_hash, role, last_login_at, created_at, updated_at
	`, email, name, avatarURL, passwordHash).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &googleID, &u.PasswordHash,
		&u.Role, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if googleID.Valid {
		u.GoogleID = googleID.String
	}
	return &u, nil
}

func (r *Repository) UpdateLastLogin(userID uuid.UUID) error {
	now := time.Now()
	_, err := r.db.Exec(`UPDATE users SET last_login_at = $1, updated_at = $1 WHERE id = $2`, now, userID)
	return err
}

func (r *Repository) CreateLearningProfile(userID uuid.UUID) error {
	_, err := r.db.Exec(`
		INSERT INTO learning_profiles (user_id)
		VALUES ($1) ON CONFLICT (user_id) DO NOTHING
	`, userID)
	return err
}
