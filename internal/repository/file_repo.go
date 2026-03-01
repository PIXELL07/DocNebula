// FYI :
// DB access layer
// This file contains the FileRepo which manages the persistence of file records associated with jobs.
// It provides methods to create new file entries, update their processing status, and retrieve files by job ID.
// The repository uses SQL queries to interact with the database and ensures proper error handling and resource management.
package repository

import (
	"DocNebula/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// FileRepo handles persistence for files extracted from a job.
type FileRepo struct {
	DB *sql.DB
}

// Create inserts a new file record for a job.
func (r *FileRepo) Create(ctx context.Context, jobID, path string) (*models.File, error) {
	id := uuid.NewString()

	q := `
	INSERT INTO files (id, job_id, path, status)
	VALUES ($1,$2,$3,$4)
	RETURNING created_at, updated_at
	`

	var createdAt, updatedAt sql.NullTime
	if err := r.DB.QueryRowContext(ctx, q,
		id,
		jobID,
		path,
		models.FilePending,
	).Scan(&createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("file create failed: %w", err)
	}

	return &models.File{
		ID:        id,
		JobID:     jobID,
		Path:      path,
		Status:    models.FilePending,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}, nil
}

// UpdateStatus safely updates the file processing state.
func (r *FileRepo) UpdateStatus(ctx context.Context, fileID string, status models.FileStatus) error {
	res, err := r.DB.ExecContext(ctx,
		`UPDATE files SET status=$1, updated_at=NOW() WHERE id=$2`,
		status,
		fileID,
	)
	if err != nil {
		return fmt.Errorf("file status update failed: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("file not found")
	}
	return nil
}

// GetByJob lists files for a job (useful for orchestration).
func (r *FileRepo) GetByJob(ctx context.Context, jobID string) ([]models.File, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, job_id, path, status, created_at, updated_at
		 FROM files WHERE job_id=$1`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.File
	for rows.Next() {
		var f models.File
		if err := rows.Scan(
			&f.ID,
			&f.JobID,
			&f.Path,
			&f.Status,
			&f.CreatedAt,
			&f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}
