package repository

import (
	"DocNebula/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// JobRepo handles persistence for jobs.
type JobRepo struct {
	DB *sql.DB
}

// Create is idempotent: same idempotency_key returns same job.
func (r *JobRepo) Create(ctx context.Context, idemKey string) (*models.Job, error) {
	// First try to find existing job
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
		// (idempotent hit)
		return &job, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("job lookup failed: %w", err)
	}

	// create new job
	id := uuid.NewString()

	q := `
	INSERT INTO jobs (id, status, retry_count, idempotency_key)
	VALUES ($1,$2,0,$3)
	RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime

	err = r.DB.QueryRowContext(ctx, q,
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
		CreatedAt:      createdAt.Time,
		UpdatedAt:      updatedAt.Time,
	}, nil
}

func (r *JobRepo) Get(ctx context.Context, id string) (*models.Job, error) {
	var job models.Job

	err := r.DB.QueryRowContext(ctx, `
		SELECT id, status, retry_count, idempotency_key, created_at, updated_at
		FROM jobs WHERE id=$1
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
		return nil, err
	}

	return &job, nil
}

// UpdateStatus updates job state safely.
func (r *JobRepo) UpdateStatus(ctx context.Context, id string, status models.JobStatus) error {
	res, err := r.DB.ExecContext(ctx,
		`UPDATE jobs SET status=$1, updated_at=NOW() WHERE id=$2`,
		status,
		id,
	)
	if err != nil {
		return fmt.Errorf("job status update failed: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("job not found")
	}
	return nil
}
