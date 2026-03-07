package environments

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/jackc/pgx/v5"
)

// Used for environment creation only to ensure we're able to connect to the database
func ValidateConnection(ctx context.Context, conn DatabaseConnection) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return shared.NewAppError(http.StatusBadRequest, shared.ErrInvalidInput, fmt.Sprintf("invalid connection string: %v", err))
	}

	if config.TLSConfig == nil && conn.SSLMode != "" {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: conn.SSLMode == "require" || conn.SSLMode == "allow",
		}
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return shared.NewAppError(http.StatusBadRequest, shared.ErrInvalidInput, fmt.Sprintf("could not connect: %v", err))
	}
	defer pgConn.Close(ctx)

	if err := pgConn.Ping(ctx); err != nil {
		return shared.NewAppError(http.StatusBadRequest, shared.ErrInvalidInput, fmt.Sprintf("database did not respond: %v", err))
	}

	return nil
}

// TODO: Add support for SSL certificates with verify-full, verify-ca, and custom CA certificates.
// For now we only support "require" mode for migrations.
func ValidateConnectionAsMigrator(ctx context.Context, conn DatabaseConnection) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return shared.NewAppError(http.StatusBadRequest, shared.ErrInvalidInput, fmt.Sprintf("invalid connection string: %v", err))
	}

	if config.TLSConfig == nil && conn.SSLMode != "" {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: conn.SSLMode == "require" || conn.SSLMode == "allow",
		}
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return shared.NewAppError(http.StatusBadRequest, shared.ErrInvalidInput, fmt.Sprintf("could not connect: %v", err))
	}
	defer pgConn.Close(ctx)

	if err := pgConn.Ping(ctx); err != nil {
		return shared.NewAppError(http.StatusBadRequest, shared.ErrInvalidInput, fmt.Sprintf("database did not respond: %v", err))
	}

	var canCreate bool
	query := `SELECT has_schema_privilege(current_user, 'public', 'CREATE')`
	err = pgConn.QueryRow(ctx, query).Scan(&canCreate)
	if err != nil {
		return shared.NewAppError(http.StatusInternalServerError, shared.ErrInternalServer, "Could not check permissions")
	}

	if !canCreate {
		return shared.NewAppError(http.StatusForbidden, shared.ErrUnauthorized, "The user can connect but does not have permission to create objects in the public schema.")
	}

	return nil
}

// An environment relies on a migrator to be able to execute a migration, that's why we're using ValidateConnectionAsMigrator
func CreateProjectEnvironment(ctx context.Context, input CreateEnvironmentRequest, masterKey []byte, repo *EnvRepo) error {
	connInfo := DatabaseConnection{
		ConnectionString: input.ConnectionUrl,
		SSLMode:          input.SSLMode,
	}

	if err := ValidateConnectionAsMigrator(ctx, connInfo); err != nil {
		return fmt.Errorf("pre-storage validation failed: %w", err)
	}

	encryptedURL, err := shared.EncryptToAes256(input.ConnectionUrl, masterKey)
	if err != nil {
		return fmt.Errorf("security breach: could not encrypt string: %w", err)
	}

	return repo.Create(ctx, input.Name, input.ProjectID, string(encryptedURL), input.SSLMode)
}
