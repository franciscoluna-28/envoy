package environments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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

// @Helper Validates a database connection as a migrator user
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

// @UseCase Creates a new environment for a project
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

	err = repo.CreateEnvironment(ctx, environment)

	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	return nil
}

// @UseCase Gets the current schema of an environment
func GetEnvironmentSchema(ctx context.Context, envID string, repo Repository, masterKey []byte) ([]SchemaColumn, error) {

	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)

	if err != nil {
		return nil, fmt.Errorf("failed to decrypt connection string: %w", err)
	}

	pgConn, err := ConnectAsMigratorHelper(ctx, decryptedURL)
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

// @UseCase Gets all environments for a project
func GetAllEnvironmentsByProjectID(ctx context.Context, projectID string, repo Repository) ([]Environment, error) {
	return repo.GetAllEnvironmentsByProjectID(ctx, projectID)
}

// @UseCase Gets an environment by its ID
func GetEnvironmentByID(ctx context.Context, environmentID string, repo Repository) (*Environment, error) {
	return repo.GetEnvironmentByID(ctx, environmentID)
}

// @UseCase Updates an environment
func UpdateEnvironment(ctx context.Context, envID string, input UpdateEnvironmentRequest, masterKey []byte, repo Repository) (*Environment, error) {
	// Get existing environment
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}

	// Validate new connection
	connInfo := DatabaseConnection{
		ConnectionString: input.ConnectionUrl,
	}

	if err := ValidateDatabaseConnectionAsMigrator(ctx, connInfo); err != nil {
		return nil, fmt.Errorf("pre-storage validation failed: %w", err)
	}

	// Encrypt new connection URL
	encryptedURL, err := EncryptToAes256(input.ConnectionUrl, masterKey)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt string: %w", err)
	}

	// Update environment fields
	env.Name = input.Name
	env.Type = input.Type
	env.ConnectionStringEncrypted = string(encryptedURL)
	env.UpdatedAt = time.Now()

	err = repo.UpdateEnvironment(ctx, *env)
	if err != nil {
		return nil, fmt.Errorf("failed to update environment: %w", err)
	}

	return env, nil
}

// @UseCase Previews the schema changes that would be applied by a migration
func PreviewEnvironmentSchemaChanges(ctx context.Context, envID string, repo Repository, sqlContent string, masterKey []byte) ([]SchemaColumn, error) {
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt connection string: %w", err)
	}

	pgConn, err := ConnectAsMigratorHelper(ctx, decryptedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pgConn.Close(ctx)

	tx, err := pgConn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// ! Never forget the rollback, we're not playing with CRUDs. This is platform engineering.
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, sqlContent)
	if err != nil {
		return nil, fmt.Errorf("SQL validation failed: %w", err)
	}

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

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query schema: %w", err)
	}
	defer rows.Close()

	var columns []SchemaColumn
	for rows.Next() {
		var col SchemaColumn
		if err := rows.Scan(&col.TableName, &col.ColumnName, &col.DataType, &col.IsNullable); err != nil {
			return nil, fmt.Errorf("failed to scan schema row: %w", err)
		}
		columns = append(columns, col)
	}

	return columns, nil
}

// @UseCase Runs a database migration in an environment
func RunDatabaseMigrationInEnvironment(ctx context.Context, envID string, repo Repository, migration CreateEnvironmentMigrationRequest, userId string, masterKey []byte, checkSumKey []byte) error {

	ctx, cancel := context.WithTimeout(ctx, 8*time.Minute) // 8 minutes is a safe timeout for most migrations
	defer cancel()

	if err := repo.FindEnvironmentMigrationByClientId(ctx, migration.ClientId); err != nil {
		// If the migration doesn't exist, that's expected - continue with creation
		if errors.Is(err, ErrEnvironmentNotFound) {
			// No existing migration found, proceed with creation
		} else {
			// Some other error occurred
			return err
		}
	}

	fmt.Printf("[DEBUG] Computing checksum for SQL content (length: %d)\n", len(migration.SQLContent))
	fmt.Printf("[DEBUG] Checksum key length: %d\n", len(checkSumKey))
	computedChecksum := GenerateSQLChecksum(migration.SQLContent, checkSumKey)
	fmt.Printf("[DEBUG] Computed checksum: %s (length: %d)\n", computedChecksum, len(computedChecksum))

	migrationID := uuid.New().String()
	now := time.Now()
	envMigration := EnvironmentMigration{
		ID:            migrationID,
		Checksum:      computedChecksum,
		EnvironmentID: migration.EnvironmentID,
		Name:          migration.Name,
		Description:   migration.Description,
		SQLContent:    migration.SQLContent,
		Status:        MigrationStatusRunning,
		ExecutedAt:    &now,
		CreatedAt:     now,
		Duration:      0,
		ExecutedBy:    userId,
		ErrorMessage:  "",
		ClientId:      migration.ClientId,
	}

	fmt.Printf("[DEBUG] Creating migration record with ID: %s\n", migrationID)
	if err := repo.CreateEnvironmentMigration(ctx, envMigration); err != nil {
		fmt.Printf("[DEBUG] Failed to create migration record: %v\n", err)
		return fmt.Errorf("failed to create migration record: %w", err)
	}

	startTime := time.Now()

	fmt.Printf("[DEBUG] Getting environment by ID: %s\n", envID)
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to get environment: %v\n", err)
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, err, startTime)
	}

	fmt.Printf("[DEBUG] Decrypting connection string (length: %d)\n", len(env.ConnectionStringEncrypted))
	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to decrypt connection string: %v\n", err)
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, fmt.Errorf("failed to decrypt: %w", err), startTime)
	}

	fmt.Printf("[DEBUG] Connecting to database...\n")
	pgConn, err := ConnectAsMigratorHelper(ctx, decryptedURL)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to connect to database: %v\n", err)
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, err, startTime)
	}
	defer pgConn.Close(ctx)

	fmt.Printf("[DEBUG] Beginning transaction...\n")
	tx, err := pgConn.Begin(ctx)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to begin transaction: %v\n", err)
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, err, startTime)
	}

	defer tx.Rollback(ctx)

	fmt.Printf("[DEBUG] Executing SQL migration (length: %d)...\n", len(migration.SQLContent))
	_, err = tx.Exec(ctx, migration.SQLContent)
	if err != nil {
		fmt.Printf("[DEBUG] SQL execution failed: %v\n", err)
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, fmt.Errorf("migration execution failed: %w", err), startTime)
	}

	fmt.Printf("[DEBUG] Committing transaction...\n")
	if err := tx.Commit(ctx); err != nil {
		fmt.Printf("[DEBUG] Failed to commit transaction: %v\n", err)
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, fmt.Errorf("failed to commit migration: %w", err), startTime)
	}

	executionTime := time.Since(startTime).Milliseconds()
	envMigration.Status = MigrationStatusCompleted
	envMigration.Duration = executionTime
	envMigration.ErrorMessage = ""

	fmt.Printf("[DEBUG] Updating migration record with completion status (duration: %dms)...\n", executionTime)
	if err := repo.UpdateEnvironmentMigration(ctx, envMigration); err != nil {
		fmt.Printf("Warning: failed to update migration record: %v\n", err)
	}

	fmt.Printf("[DEBUG] Migration completed successfully!\n")
	return nil
}

// @UseCase Verifies that the database user has the least privilege required
func VerifyDatabasePermissions(ctx context.Context, envID string, repo Repository, masterKey []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	config, err := pgx.ParseConfig(decryptedURL)
	if err != nil {
		return err
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer pgConn.Close(ctx)

	// TODO: Even if it's a bad practice, create a flag to skip this check (I won't be responsible for any damage caused by this suggestion)
	var canCreate bool
	err = pgConn.QueryRow(ctx, `SELECT has_schema_privilege(current_user, 'public', 'CREATE')`).Scan(&canCreate)
	if err != nil {
		return err
	}

	if canCreate {
		return fmt.Errorf("security breach: user can MODIFY schema (DDL allowed)")
	}

	var canOperateData bool

	// TODO: Consider making this customizable based on the environment's requirements in the future
	queryDML := `
        SELECT bool_and(has_table_privilege(current_user, table_name, 'SELECT, INSERT, UPDATE, DELETE'))
        FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_type = 'BASE TABLE';`

	err = pgConn.QueryRow(ctx, queryDML).Scan(&canOperateData)
	if err != nil {
		return err
	}

	if !canOperateData {
		return fmt.Errorf("insufficient privileges: user cannot perform basic CRUD operations")
	}

	return nil
}

// @UseCase Gets all migrations for an environment
func GetAllMigrationsByEnvironmentID(ctx context.Context, envID string, repo Repository) ([]EnvironmentMigration, error) {
	return repo.GetEnvironmentMigrations(ctx, envID)
}

// @UseCase Gets a migration by ID
func GetEnvironmentMigrationByID(ctx context.Context, migrationID string, repo Repository) (*EnvironmentMigration, error) {
	return repo.GetEnvironmentMigrationByID(ctx, migrationID)
}

func ValidateEnvironmentConnection(ctx context.Context, envID string, repo Repository, masterKey []byte) error {
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	return ValidateDatabaseConnectionAsMigrator(ctx, DatabaseConnection{
		ConnectionString: decryptedURL,
	})
}
