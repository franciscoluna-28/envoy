package environments

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateEnvironment(ctx context.Context, env Environment) error
	GetAllEnvironmentsByProjectID(ctx context.Context, projectID string) ([]Environment, error)
	GetEnvironmentByID(ctx context.Context, id string) (*Environment, error)
	UpdateEnvironment(ctx context.Context, env Environment) error
	DeleteEnvironment(ctx context.Context, id string) error
	CreateMigration(ctx context.Context, envID string, sqlContent string) error
	CreateEnvironmentDbUser(ctx context.Context, user EnvironmentDbUser) error
	GetEnvironmentDbUsers(ctx context.Context, envID string) ([]EnvironmentDbUser, error)
	CreateEnvironmentMigration(ctx context.Context, migration EnvironmentMigration) error
	GetEnvironmentMigrations(ctx context.Context, envID string) ([]EnvironmentMigration, error)
	GetEnvironmentMigrationByID(ctx context.Context, id string) (*EnvironmentMigration, error)
	UpdateEnvironmentMigration(ctx context.Context, migration EnvironmentMigration) error
	FindEnvironmentMigrationByClientId(ctx context.Context, clientId string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateEnvironment(ctx context.Context, env Environment) error {
	query := `INSERT INTO environments (id, name, project_id, type, connection_string_encrypted, connection_error, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			  `

	_, err := r.db.ExecContext(ctx, query, env.ID, env.Name, env.ProjectID, env.Type, env.ConnectionStringEncrypted, env.ConnectionError, env.CreatedAt, env.UpdatedAt)
	if err != nil {
		return err
	}
	return err
}

func (r *repository) GetAllEnvironmentsByProjectID(ctx context.Context, projectID string) ([]Environment, error) {
	query := `SELECT id, name, project_id, type, connection_string_encrypted, connection_error, created_at, updated_at 
			  FROM environments WHERE project_id = ?`

	var envs []Environment

	err := r.db.SelectContext(ctx, &envs, query, projectID)

	if err != nil {
		return envs, err
	}

	return envs, err
}

func (r *repository) GetEnvironmentByID(ctx context.Context, id string) (*Environment, error) {
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

func (r *repository) UpdateEnvironment(ctx context.Context, env Environment) error {
	query := `UPDATE environments SET name = ?, project_id = ?, type = ?, connection_string_encrypted = ?, connection_error = ?, updated_at = ? 
			  WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, env.Name, env.ProjectID, env.Type, env.ConnectionStringEncrypted, env.ConnectionError, env.UpdatedAt, env.ID)
	return err
}

func (r *repository) DeleteEnvironment(ctx context.Context, id string) error {
	query := `DELETE FROM environments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) CreateMigration(ctx context.Context, envID string, sqlContent string) error {
	query := `INSERT INTO migrations (env_id, sql_content) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, envID, sqlContent)
	return err
}

func (r *repository) CreateEnvironmentDbUser(ctx context.Context, user EnvironmentDbUser) error {
	query := `INSERT INTO environment_db_users (id, environment_id, name, connection_string_encrypted, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.EnvironmentID, user.Name, user.ConnectionStringEncrypted, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *repository) GetEnvironmentDbUsers(ctx context.Context, envID string) ([]EnvironmentDbUser, error) {
	query := `SELECT id, environment_id, name, connection_string_encrypted, created_at, updated_at 
			  FROM environment_db_users WHERE environment_id = ?`
	var users []EnvironmentDbUser
	err := r.db.SelectContext(ctx, &users, query, envID)
	return users, err
}

func (r *repository) CreateEnvironmentMigration(ctx context.Context, migration EnvironmentMigration) error {
	query := `INSERT INTO environment_migrations (id, checksum, environment_id, name, description, sql_content, status, executed_at, created_at, duration, executed_by, error_message, client_id) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		migration.ID, migration.Checksum, migration.EnvironmentID, migration.Name, migration.Description,
		migration.SQLContent, migration.Status, migration.ExecutedAt, migration.CreatedAt,
		migration.Duration, migration.ExecutedBy, migration.ErrorMessage, migration.ClientId)
	return err
}

func (r *repository) GetEnvironmentMigrations(ctx context.Context, envID string) ([]EnvironmentMigration, error) {
	query := `SELECT id, checksum, environment_id, name, description, sql_content, status, executed_at, created_at, duration, executed_by, error_message, client_id 
			  FROM environment_migrations WHERE environment_id = ? ORDER BY created_at DESC`
	var migrations []EnvironmentMigration
	err := r.db.SelectContext(ctx, &migrations, query, envID)
	return migrations, err
}

func (r *repository) GetEnvironmentMigrationByID(ctx context.Context, id string) (*EnvironmentMigration, error) {
	query := `SELECT id, checksum, environment_id, name, description, sql_content, status, executed_at, created_at, duration, executed_by, error_message, client_id 
			  FROM environment_migrations WHERE id = ?`
	var migration EnvironmentMigration
	err := r.db.GetContext(ctx, &migration, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEnvironmentNotFound
		}
		return nil, err
	}
	return &migration, nil
}

func (r *repository) UpdateEnvironmentMigration(ctx context.Context, migration EnvironmentMigration) error {
	query := `UPDATE environment_migrations SET 
			  name = ?, description = ?, sql_content = ?, status = ?, executed_at = ?, 
			  duration = ?, executed_by = ?, error_message = ?, client_id = ?
			  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		migration.Name, migration.Description, migration.SQLContent, migration.Status, migration.ExecutedAt,
		migration.Duration, migration.ExecutedBy, migration.ErrorMessage, migration.ClientId, migration.ID)
	return err
}

func (r *repository) FindEnvironmentMigrationByClientId(ctx context.Context, migrationId string) error {
	query := `SELECT id, environment_id, name, description, sql_content, status, executed_at, created_at, duration, executed_by, error_message, client_id 
			  FROM environment_migrations WHERE client_id = ?`
	var migration EnvironmentMigration
	err := r.db.GetContext(ctx, &migration, query, migrationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEnvironmentNotFound
		}
		return err
	}
	return nil
}
