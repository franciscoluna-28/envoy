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

// @UseCase Creates a new environment for a project
func CreateProjectEnvironment(ctx context.Context, input CreateEnvironmentRequest, masterKey []byte, repo Repository) error {
	connInfo := DatabaseConnection{
		ConnectionString: input.ConnectionUrl,
	}

	dbVersion, err := ValidateDatabaseConnectionAsMigrator(ctx, connInfo)

	if err != nil {
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
		DbEngine:                  dbVersion,
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

	_, err = ValidateDatabaseConnectionAsMigrator(ctx, connInfo)
	if err != nil {
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

	if err := repo.FindEnvironmentMigrationByClientId(ctx, migration.ClientId); err != nil && !errors.Is(err, ErrEnvironmentNotFound) {
		return err
	}

	computedChecksum := GenerateSQLChecksum(migration.SQLContent, checkSumKey)

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

	if err := repo.CreateEnvironmentMigration(ctx, envMigration); err != nil {
		return fmt.Errorf("failed to create migration record: %w", err)
	}

	startTime := time.Now()

	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, err, startTime)
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, fmt.Errorf("failed to decrypt: %w", err), startTime)
	}

	pgConn, err := ConnectAsMigratorHelper(ctx, decryptedURL)
	if err != nil {
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, err, startTime)
	}
	defer pgConn.Close(ctx)

	tx, err := pgConn.Begin(ctx)
	if err != nil {
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, err, startTime)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, migration.SQLContent)
	if err != nil {
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, fmt.Errorf("migration execution failed: %w", err), startTime)
	}

	if err := tx.Commit(ctx); err != nil {
		return UpdateMigrationErrorHelper(ctx, repo, migrationID, fmt.Errorf("failed to commit migration: %w", err), startTime)
	}

	executionTime := time.Since(startTime).Milliseconds()
	envMigration.Status = MigrationStatusCompleted
	envMigration.Duration = executionTime
	envMigration.ErrorMessage = ""

	if err := repo.UpdateEnvironmentMigration(ctx, envMigration); err != nil {
		return fmt.Errorf("failed to update migration record: %w", err)
	}

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

// @UseCase Validates that an environment can connect to the database as a migrator
func ValidateEnvironmentConnection(ctx context.Context, envID string, repo Repository, masterKey []byte) error {
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	_, err = ValidateDatabaseConnectionAsMigrator(ctx, DatabaseConnection{
		ConnectionString: decryptedURL,
	})
	return err
}

// @UseCase Audits permissions before running a migration (preview schema mode)
func AuditPermissionsBeforeMigration(ctx context.Context, envID string, repo Repository, sqlContent string, targetUser string, masterKey []byte) ([]TablePermission, error) {
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return nil, err
	}

	config, _ := pgx.ParseConfig(decryptedURL)
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback, this is not an actual migration

	// Only execute the migration in the transaction context
	if _, err := tx.Exec(ctx, sqlContent); err != nil {
		return nil, fmt.Errorf("migration preview failed: %w", err)
	}

	// Change the app user to verify the new permissions we will give
	setRole := fmt.Sprintf("SET ROLE %s", pgx.Identifier{targetUser}.Sanitize())
	if _, err := tx.Exec(ctx, setRole); err != nil {
		return nil, fmt.Errorf("app user '%s' not found in DB: %w", targetUser, err)
	}

	query := `
        SELECT table_name, privilege_type 
        FROM information_schema.table_privileges 
        WHERE table_schema = 'public' 
        AND grantee = current_user;`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permMap := make(map[string][]string)
	for rows.Next() {
		var table, priv string
		rows.Scan(&table, &priv)
		permMap[table] = append(permMap[table], priv)
	}

	var report []TablePermission
	for table, privs := range permMap {
		report = append(report, TablePermission{
			TableName:  table,
			Privileges: privs,
			IsMissing:  len(privs) < 2, // Missing if they cannot do at least 2 operations (SELECT and INSERT)
		})
	}

	return report, nil
}

// @UseCase Audits the current permissions of a specific database role
func AuditPermissionsWithCurrentSchema(ctx context.Context, envID string, repo Repository, dbUser string, masterKey []byte) ([]TablePermission, error) {
	env, err := repo.GetEnvironmentByID(ctx, envID)
	if err != nil {
		return nil, err
	}

	decryptedURL, err := DecryptFromAes256(env.ConnectionStringEncrypted, masterKey)
	if err != nil {
		return nil, err
	}

	config, err := pgx.ParseConfig(decryptedURL)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}
	defer conn.Close(ctx)

	// Transaction only to isolate the SET ROLE command
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	setRole := fmt.Sprintf("SET ROLE %s", pgx.Identifier{dbUser}.Sanitize())
	if _, err := tx.Exec(ctx, setRole); err != nil {
		return nil, fmt.Errorf("app user '%s' not found in DB: %w", dbUser, err)
	}

	query := `
        SELECT table_name, privilege_type 
        FROM information_schema.table_privileges 
        WHERE table_schema = 'public' 
        AND grantee = current_user;`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permMap := make(map[string][]string)
	for rows.Next() {
		var table, priv string
		if err := rows.Scan(&table, &priv); err != nil {
			return nil, err
		}
		permMap[table] = append(permMap[table], priv)
	}

	var report []TablePermission
	for table, privs := range permMap {
		report = append(report, TablePermission{
			TableName:  table,
			Privileges: privs,
			IsMissing:  len(privs) < 2, // Missing if they cannot do at least 2 operations (SELECT and INSERT)
		})
	}

	return report, nil
}
