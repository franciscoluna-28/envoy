package repositories

import (
	"context"
	"server/internal/api/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, created_at)
        VALUES (:id, :email, :password_hash, :created_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = ?`

	err := r.db.GetContext(ctx, &user, query, string(id))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`

	err := r.db.GetContext(ctx, &user, query, email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
