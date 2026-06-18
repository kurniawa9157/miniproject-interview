package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"jumpapay/backend/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// UpsertByGoogleID inserts or updates a user from Google OAuth data.
func (r *UserRepository) UpsertByGoogleID(ctx context.Context, googleID, name, email, photoURL string) (*model.User, error) {
	query := `
		INSERT INTO users (google_id, name, email, photo_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id) DO UPDATE
		SET name = EXCLUDED.name,
		    email = EXCLUDED.email,
		    photo_url = EXCLUDED.photo_url
		RETURNING id, google_id, name, email, COALESCE(photo_url, ''), is_admin, created_at
	`

	var u model.User
	err := r.db.QueryRow(ctx, query, googleID, name, email, photoURL).Scan(
		&u.ID, &u.GoogleID, &u.Name, &u.Email, &u.PhotoURL, &u.IsAdmin, &u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert user: %w", err)
	}

	return &u, nil
}

// FindByID returns a user by UUID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, google_id, name, email, COALESCE(photo_url, ''), is_admin, created_at
		FROM users WHERE id = $1
	`

	var u model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.GoogleID, &u.Name, &u.Email, &u.PhotoURL, &u.IsAdmin, &u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &u, nil
}
