package repository

import (
	"DocNebula/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobRepo struct {
	DB *sql.DB
}

// Create is idempotent: same idempotency_key returns same job.
func (r *JobRepo) Create(ctx context.Context, idemKey string) (*models.Job, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var job models.Job

	row := r.DB.QueryRowContext(ctx, `
		SELECT id, status, retry_count, idempotency_key, created_at, updated_at
		FROM jobs
		WHERE idempotency_key=$1
	`, idemKey)

	err := row.Scan(
		&job.ID,
		&job.Status,
		&job.RetryCount,
		&job.IdempotencyKey,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err == nil {
		return &job, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("job lookup failed: %w", err)
	}

	id := uuid.NewString()

	var createdAt, updatedAt time.Time

	err = r.DB.QueryRowContext(ctx, `
		INSERT INTO jobs (id, status, retry_count, idempotency_key)
		VALUES ($1,$2,0,$3)
		RETURNING created_at, updated_at
	`,
		id,
		models.JobUploaded,
		idemKey,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("job insert failed: %w", err)
	}

	return &models.Job{
		ID:             id,
		Status:         models.JobUploaded,
		RetryCount:     0,
		IdempotencyKey: idemKey,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}

func (r *JobRepo) Get(ctx context.Context, id string) (*models.Job, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var job models.Job

	err := r.DB.QueryRowContext(ctx, `
		SELECT id, status, retry_count, idempotency_key, created_at, updated_at
		FROM jobs
		WHERE id=$1
	`, id).Scan(
		&job.ID,
		&job.Status,
		&job.RetryCount,
		&job.IdempotencyKey,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get job failed: %w", err)
	}

	return &job, nil
}

func (r *JobRepo) UpdateStatus(ctx context.Context, id string, status models.JobStatus) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.DB.ExecContext(ctx,
		`UPDATE jobs SET status=$1, updated_at=NOW() WHERE id=$2`,
		status,
		id,
	)

	if err != nil {
		return fmt.Errorf("job status update failed: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}

	if rows == 0 {
		return errors.New("job not found")
	}

	return nil
}
