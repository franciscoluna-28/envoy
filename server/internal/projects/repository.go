package projects

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, project Project) error
	GetByID(ctx context.Context, id string, userID string) (*Project, error)
	GetAllByUserID(ctx context.Context, userID string) ([]Project, error)
	Update(ctx context.Context, project Project, userID string) error
	Delete(ctx context.Context, id string, userID string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, project Project) error {
	query := `INSERT INTO projects (id, name, created_by, created_at, updated_at) VALUES (:id, :name, :created_by, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, project)
	return err
}

func (r *repository) GetByID(ctx context.Context, id string, userID string) (*Project, error) {
	var project Project
	query := `SELECT id, name, created_by, created_at, updated_at FROM projects WHERE id = ? AND created_by = ?`
	err := r.db.GetContext(ctx, &project, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

func (r *repository) GetAllByUserID(ctx context.Context, userID string) ([]Project, error) {
	var projects []Project
	query := `SELECT id, name, created_by, created_at, updated_at FROM projects WHERE created_by = ?`
	err := r.db.SelectContext(ctx, &projects, query, userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *repository) Update(ctx context.Context, project Project, userID string) error {
	query := `UPDATE projects SET name = :name, updated_at = :updated_at WHERE id = :id AND created_by = :created_by`
	_, err := r.db.NamedExecContext(ctx, query, project)
	return err
}

func (r *repository) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM projects WHERE id = ? AND created_by = ?`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}
