// DB access layer

package repository

import (
	"DocNebula/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type JobRepository struct {
	DB *sql.DB
}

// Idempotent job creation
func (r *JobRepository) CreateJob(ctx context.Context, idemKey string) (*models.Job, error) {
	// Check if already exists
	var existingID string
	err := r.DB.QueryRowContext(ctx,
		`SELECT id FROM jobs WHERE idempotency_key=$1`,
		idemKey,
	).Scan(&existingID)

	if err == nil {
		return &models.Job{ID: existingID}, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	id := uuid.NewString()

	_, err = r.DB.ExecContext(ctx, `
		INSERT INTO jobs(id, status, idempotency_key)
		VALUES ($1,$2,$3)
	`, id, models.JobUploaded, idemKey)
	if err != nil {
		return nil, err
	}

	return &models.Job{ID: id, Status: models.JobUploaded}, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id string, status models.JobStatus) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE jobs SET status=$1, updated_at=NOW() WHERE id=$2`,
		status, id,
	)
	return err
}

func (r *JobRepository) IncrementRetry(ctx context.Context, id string) (int, error) {
	var retry int
	err := r.DB.QueryRowContext(ctx, `
		UPDATE jobs
		SET retry_count = retry_count + 1
		WHERE id=$1
		RETURNING retry_count
	`, id).Scan(&retry)
	return retry, err
}
