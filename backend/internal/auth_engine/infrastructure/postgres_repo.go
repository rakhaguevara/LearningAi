package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, provider, provider_id, name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Provider, user.ProviderID,
		user.Name, user.AvatarURL, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *PostgresAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, provider, provider_id, name, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresAuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, provider, provider_id, name, avatar_url, created_at, updated_at
		FROM users WHERE email = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *PostgresAuthRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, name = $3, avatar_url = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.Name, user.AvatarURL, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *PostgresAuthRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	return err
}

func (r *PostgresAuthRepository) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var rt domain.RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	return &rt, nil
}

func (r *PostgresAuthRepository) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, id)
	return err
}

func (r *PostgresAuthRepository) DeleteRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

func (r *PostgresAuthRepository) CreateLoginHistory(ctx context.Context, history *domain.LoginHistory) error {
	query := `
		INSERT INTO login_history (id, user_id, ip, device, location, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		history.ID, history.UserID, history.IP, history.Device, history.Location, history.CreatedAt,
	)
	return err
}

func (r *PostgresAuthRepository) scanUser(row *sql.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Provider, &u.ProviderID,
		&u.Name, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}
