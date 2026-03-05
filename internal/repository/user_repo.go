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

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id := uuid.NewString()

	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO users (id,email,password_hash)
		 VALUES ($1,$2,$3)`,
		id,
		email,
		passwordHash,
	)

	if err != nil {
		return nil, fmt.Errorf("user create failed: %w", err)
	}

	return &models.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx,
		`SELECT id,email,password_hash,created_at
		 FROM users
		 WHERE email=$1`,
		email,
	)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("user lookup failed: %w", err)
	}

	return &u, nil
}
