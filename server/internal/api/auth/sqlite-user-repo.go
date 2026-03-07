package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `
        INSERT INTO users (id, email, password_hash, created_at)
        VALUES (:id, :email, :password_hash, :created_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id UserID) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = ?`

	err := r.db.GetContext(ctx, &user, query, string(id))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`

	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
