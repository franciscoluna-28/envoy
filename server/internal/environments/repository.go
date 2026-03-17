package environments

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, env Environment) error
	GetAllByProjectID(ctx context.Context, projectID string) ([]Environment, error)
	GetByID(ctx context.Context, id string) (*Environment, error)
	Update(ctx context.Context, env Environment) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, env Environment) error {
	query := `INSERT INTO environments (id, name, project_id, type, connection_string_encrypted, connection_error, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			  `

	_, err := r.db.ExecContext(ctx, query, env.ID, env.Name, env.ProjectID, env.Type, env.ConnectionStringEncrypted, env.ConnectionError, env.CreatedAt, env.UpdatedAt)
	if err != nil {
		return err
	}
	return err
}

func (r *repository) GetAllByProjectID(ctx context.Context, projectID string) ([]Environment, error) {
	query := `SELECT id, name, project_id, type, connection_string_encrypted, connection_error, created_at, updated_at 
			  FROM environments WHERE project_id = ?`

	var envs []Environment

	err := r.db.SelectContext(ctx, &envs, query, projectID)

	if err != nil {
		return envs, err
	}

	return envs, err
}

func (r *repository) GetByID(ctx context.Context, id string) (*Environment, error) {
	query := `SELECT id, name, project_id, type, connection_string_encrypted, connection_error, created_at, updated_at 
			  FROM environments WHERE id = ?`

	var env Environment
	err := r.db.GetContext(ctx, &env, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}

	return &env, nil
}

func (r *repository) Update(ctx context.Context, env Environment) error {
	query := `UPDATE environments SET name = ?, project_id = ?, type = ?, connection_string_encrypted = ?, connection_error = ?, updated_at = ? 
			  WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, env.Name, env.ProjectID, env.Type, env.ConnectionStringEncrypted, env.ConnectionError, env.UpdatedAt, env.ID)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM environments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
