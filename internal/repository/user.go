package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourusername/resume-builder/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (github_id, username, email, name, avatar_url, encrypted_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	now := time.Now()
	return r.db.QueryRowContext(
		ctx, query,
		user.GitHubID, user.Username, user.Email, user.Name, user.AvatarURL, user.EncryptedToken, now, now,
	).Scan(&user.ID)
}

func (r *UserRepository) GetByGitHubID(ctx context.Context, githubID int64) (*model.User, error) {
	query := `
		SELECT id, github_id, username, email, name, avatar_url, encrypted_token, created_at, updated_at
		FROM users
		WHERE github_id = $1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, githubID).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.Name,
		&user.AvatarURL, &user.EncryptedToken, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, github_id, username, email, name, avatar_url, encrypted_token, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.Email, &user.Name,
		&user.AvatarURL, &user.EncryptedToken, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateToken(ctx context.Context, userID int64, encryptedToken string) error {
	query := `UPDATE users SET encrypted_token = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, encryptedToken, time.Now(), userID)
	return err
}
