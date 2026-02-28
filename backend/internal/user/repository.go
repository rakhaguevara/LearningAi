package user

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

func (r *Repository) FindByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT id, email, name, avatar_url, google_id, role, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.GoogleID,
		&u.Role, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetLearningProfile(userID uuid.UUID) (*models.LearningProfile, error) {
	var p models.LearningProfile
	err := r.db.QueryRow(`
		SELECT id, user_id, preferred_style, difficulty_level, weekly_target_hours, created_at, updated_at
		FROM learning_profiles WHERE user_id = $1
	`, userID).Scan(
		&p.ID, &p.UserID, &p.PreferredStyle, &p.DifficultyLevel,
		&p.WeeklyTargetHours, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) GetInterestTags(userID uuid.UUID) ([]models.InterestTag, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, tag, category, weight, created_at
		FROM interest_tags WHERE user_id = $1 ORDER BY weight DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.InterestTag
	for rows.Next() {
		var t models.InterestTag
		if err := rows.Scan(&t.ID, &t.UserID, &t.Tag, &t.Category, &t.Weight, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *Repository) AddInterestTag(userID uuid.UUID, tag, category string, weight float64) error {
	_, err := r.db.Exec(`
		INSERT INTO interest_tags (user_id, tag, category, weight)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, tag) DO UPDATE SET weight = $4, category = $3
	`, userID, tag, category, weight)
	return err
}
