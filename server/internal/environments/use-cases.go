package environments

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func ValidateDatabaseConnection(ctx context.Context, conn DatabaseConnection) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return fmt.Errorf("invalid connection string: %w", err)
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}
	defer pgConn.Close(ctx)

	if err := pgConn.Ping(ctx); err != nil {
		return fmt.Errorf("database did not respond: %w", err)
	}

	return nil
}

func ValidateDatabaseConnectionAsMigrator(ctx context.Context, conn DatabaseConnection) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return fmt.Errorf("invalid connection string: %w", err)
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}
	defer pgConn.Close(ctx)

	if err := pgConn.Ping(ctx); err != nil {
		return fmt.Errorf("database did not respond: %w", err)
	}

	var canCreate bool
	query := `SELECT has_schema_privilege(current_user, 'public', 'CREATE')`
	err = pgConn.QueryRow(ctx, query).Scan(&canCreate)
	if err != nil {
		return fmt.Errorf("could not check permissions: %w", err)
	}

	if !canCreate {
		return fmt.Errorf("the user can connect but does not have permission to create objects in the public schema")
	}

	return nil
}

func CreateProjectEnvironment(ctx context.Context, input CreateEnvironmentRequest, masterKey []byte, repo Repository) error {
	connInfo := DatabaseConnection{
		ConnectionString: input.ConnectionUrl,
	}

	if err := ValidateDatabaseConnectionAsMigrator(ctx, connInfo); err != nil {
		return fmt.Errorf("pre-storage validation failed: %w", err)
	}

	encryptedURL, err := EncryptToAes256(input.ConnectionUrl, masterKey)

	if err != nil {
		return fmt.Errorf("could not encrypt string: %w", err)
	}

	now := time.Now()
	var environment Environment = Environment{
		ID:                        uuid.New().String(),
		Type:                      input.Type,
		Name:                      input.Name,
		ProjectID:                 input.ProjectID,
		ConnectionStringEncrypted: string(encryptedURL),
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	err = repo.Create(ctx, environment)

	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	return nil
}

func GetEnvironmentSchema(ctx context.Context, envID string, repo Repository, masterKey []byte) ([]SchemaColumn, error) {

	env, err := repo.GetByID(ctx, envID)
	if err != nil {
		return nil, err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)

	if err != nil {
		return nil, fmt.Errorf("failed to decrypt connection string: %w", err)
	}

	config, err := pgx.ParseConfig(decryptedURL)

	if err != nil {
		return nil, err
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Printf("GetEnvironmentSchema: Database connection failed - %v\n", err)
		return nil, err
	}
	defer pgConn.Close(ctx)

	// PostgreSQL schema introspection query
	query := `
        SELECT 
            table_name, 
            column_name, 
            data_type, 
            is_nullable 
        FROM 
            information_schema.columns 
        WHERE 
            table_schema = 'public'
        ORDER BY 
            table_name, ordinal_position;`

	rows, err := pgConn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying schema failed: %w", err)
	}
	defer rows.Close()

	var columns []SchemaColumn
	for rows.Next() {
		var col SchemaColumn
		if err := rows.Scan(&col.TableName, &col.ColumnName, &col.DataType, &col.IsNullable); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func GetAllEnvironmentsByProjectID(ctx context.Context, projectID string, repo Repository) ([]Environment, error) {
	return repo.GetAllByProjectID(ctx, projectID)
}

func GetEnvironmentByID(ctx context.Context, environmentID string, repo Repository) (*Environment, error) {
	return repo.GetByID(ctx, environmentID)
}
