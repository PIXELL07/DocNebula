// how it works:
// 1. When a file is created, we bulk-insert page records with "done=false" and page numbers.
// 2. As each page is processed, we update the corresponding record to "done=true" and store the extracted text.
// 3. The NextPending method allows us to fetch the next unprocessed page, enabling resume support if the worker restarts.
// ** This design allows us to track progress at the page level and easily resume processing without losing state.

package repository

import (
	"DocNebula/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// PageRepo tracks per-page OCR progress
type PageRepo struct {
	DB *sql.DB
}

// CreatePages bulk-creates page checkpoints for a file.
func (r *PageRepo) CreatePages(ctx context.Context, fileID string, totalPages int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO pages (id, file_id, page_num, done)
		 VALUES ($1,$2,$3,false)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 1; i <= totalPages; i++ {
		if _, err := stmt.ExecContext(ctx,
			uuid.NewString(),
			fileID,
			i,
		); err != nil {
			return fmt.Errorf("page insert failed: %w", err)
		}
	}

	return tx.Commit()
}

// MarkDone marks a single page as processed.
func (r *PageRepo) MarkDone(ctx context.Context, fileID string, pageNum int, text string) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE pages
		SET done=true, text=$1
		WHERE file_id=$2 AND page_num=$3
	`, text, fileID, pageNum)
	if err != nil {
		return fmt.Errorf("page mark done failed: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("page not found")
	}
	return nil
}

// NextPending returns the next unprocessed page (resume support).
func (r *PageRepo) NextPending(ctx context.Context, fileID string) (*models.Page, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, file_id, page_num, text, done, created_at, updated_at
		FROM pages
		WHERE file_id=$1 AND done=false
		ORDER BY page_num ASC
		LIMIT 1
	`, fileID)

	var p models.Page
	if err := row.Scan(
		&p.ID,
		&p.FileID,
		&p.PageNum,
		&p.Text,
		&p.Done,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}
