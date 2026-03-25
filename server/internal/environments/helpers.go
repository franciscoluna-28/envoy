package environments

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// @Helper Connects to the database as a migrator user
func ConnectAsMigratorHelper(ctx context.Context, connectionURL string) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	pgConn, err := pgx.Connect(ctx, connectionURL)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	var canCreate bool
	query := `SELECT has_schema_privilege(current_user, 'public', 'CREATE')`
	err = pgConn.QueryRow(ctx, query).Scan(&canCreate)

	if err != nil {
		bgCtx := context.Background()
		pgConn.Close(bgCtx)
		return nil, fmt.Errorf("permission check failed: %w", err)
	}

	if !canCreate {
		bgCtx := context.Background()
		pgConn.Close(bgCtx)
		return nil, fmt.Errorf("insufficient privileges: user cannot CREATE in public schema")
	}

	return pgConn, nil
}

// @Helper Updates a migration error
func UpdateMigrationErrorHelper(ctx context.Context, repo Repository, migrationID string, err error, startTime time.Time) error {
	executionTime := time.Since(startTime).Milliseconds()

	migration, getErr := repo.GetEnvironmentMigrationByID(ctx, migrationID)
	if getErr != nil {
		return fmt.Errorf("migration failed and could not update record: %v (original error: %w)", getErr, err)
	}

	migration.Status = MigrationStatusFailed
	migration.Duration = executionTime
	migration.ErrorMessage = err.Error()

	if updateErr := repo.UpdateEnvironmentMigration(ctx, *migration); updateErr != nil {
		return fmt.Errorf("migration failed and could not update record: %v (original error: %w)", updateErr, err)
	}

	return err
}

// @Helper Validates a database connection as a migrator user and returns version
func ValidateDatabaseConnectionAsMigrator(ctx context.Context, conn DatabaseConnection) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return "", fmt.Errorf("invalid connection string: %w", err)
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return "", fmt.Errorf("could not connect: %w", err)
	}
	defer pgConn.Close(ctx)

	var dbVersion string
	err = pgConn.QueryRow(ctx, "SELECT version()").Scan(&dbVersion)
	if err != nil {
		return "", fmt.Errorf("could not get database version: %w", err)
	}

	var canCreate bool
	query := `SELECT has_schema_privilege(current_user, 'public', 'CREATE')`
	err = pgConn.QueryRow(ctx, query).Scan(&canCreate)
	if err != nil {
		return "", fmt.Errorf("could not check permissions: %w", err)
	}

	if !canCreate {
		return "", fmt.Errorf("user lacks CREATE permission in public schema")
	}

	return dbVersion, nil
}
