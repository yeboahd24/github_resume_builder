package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/yourusername/resume-builder/internal/model"
)

type ResumeRepository struct {
	db *sql.DB
}

func NewResumeRepository(db *sql.DB) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (r *ResumeRepository) Create(ctx context.Context, resume *model.Resume) error {
	projectsJSON, err := json.Marshal(resume.Projects)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO resumes (user_id, title, target_role, summary, projects, skills, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	now := time.Now()
	return r.db.QueryRowContext(
		ctx, query,
		resume.UserID, resume.Title, resume.TargetRole, resume.Summary,
		projectsJSON, pq.Array(resume.Skills), resume.IsDefault, now, now,
	).Scan(&resume.ID)
}

func (r *ResumeRepository) GetByID(ctx context.Context, id int64) (*model.Resume, error) {
	query := `
		SELECT id, user_id, title, target_role, summary, projects, skills, is_default, created_at, updated_at
		FROM resumes
		WHERE id = $1`

	resume := &model.Resume{}
	var projectsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&resume.ID, &resume.UserID, &resume.Title, &resume.TargetRole, &resume.Summary,
		&projectsJSON, pq.Array(&resume.Skills), &resume.IsDefault, &resume.CreatedAt, &resume.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(projectsJSON, &resume.Projects); err != nil {
		return nil, err
	}

	return resume, nil
}

func (r *ResumeRepository) ListByUserID(ctx context.Context, userID int64) ([]model.Resume, error) {
	query := `
		SELECT id, user_id, title, target_role, summary, projects, skills, is_default, created_at, updated_at
		FROM resumes
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resumes []model.Resume
	for rows.Next() {
		var resume model.Resume
		var projectsJSON []byte

		err := rows.Scan(
			&resume.ID, &resume.UserID, &resume.Title, &resume.TargetRole, &resume.Summary,
			&projectsJSON, pq.Array(&resume.Skills), &resume.IsDefault, &resume.CreatedAt, &resume.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(projectsJSON, &resume.Projects); err != nil {
			return nil, err
		}

		resumes = append(resumes, resume)
	}

	return resumes, rows.Err()
}

func (r *ResumeRepository) Update(ctx context.Context, resume *model.Resume) error {
	projectsJSON, err := json.Marshal(resume.Projects)
	if err != nil {
		return err
	}

	query := `
		UPDATE resumes
		SET title = $1, target_role = $2, summary = $3, projects = $4, skills = $5, is_default = $6, updated_at = $7
		WHERE id = $8`

	_, err = r.db.ExecContext(
		ctx, query,
		resume.Title, resume.TargetRole, resume.Summary,
		projectsJSON, pq.Array(resume.Skills), resume.IsDefault, time.Now(), resume.ID,
	)
	return err
}

func (r *ResumeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM resumes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
