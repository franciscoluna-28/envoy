package infra

import (
	"context"

	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) auth.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *auth.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, created_at)
        VALUES (:id, :email, :password_hash, :created_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id auth.UserID) (*auth.User, error) {
	var user auth.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = ?`

	err := r.db.GetContext(ctx, &user, query, string(id))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*auth.User, error) {
	var user auth.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`

	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
