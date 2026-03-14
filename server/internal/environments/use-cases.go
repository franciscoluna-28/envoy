package environments

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/jackc/pgx/v5"
)

func EncryptToAes256(plaintext string, masterKey []byte) (string, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func ValidateDatabase(ctx context.Context, conn DatabaseConnection) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return fmt.Errorf("invalid connection string: %w", err)
	}

	if config.TLSConfig == nil && conn.SSLMode != "" {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: conn.SSLMode == "require" || conn.SSLMode == "allow",
		}
	}

	pgConn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {

	}
	defer pgConn.Close(ctx)

	if err := pgConn.Ping(ctx); err != nil {
		return fmt.Errorf("database did not respond: %w", err)
	}

	return nil
}

func ValidateDatabaseAsMigrator(ctx context.Context, conn DatabaseConnection) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(conn.ConnectionString)
	if err != nil {
		return fmt.Errorf("invalid connection string: %w", err)
	}

	if config.TLSConfig == nil && conn.SSLMode != "" {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: conn.SSLMode == "require" || conn.SSLMode == "allow",
		}
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
		SSLMode:          input.SSLMode,
	}

	if err := ValidateDatabaseAsMigrator(ctx, connInfo); err != nil {
		return fmt.Errorf("pre-storage validation failed: %w", err)
	}

	encryptedURL, err := EncryptToAes256(input.ConnectionUrl, masterKey)
	if err != nil {
		return fmt.Errorf("security breach: could not encrypt string: %w", err)
	}

	var environment Environment = Environment{
		Name:                      input.Name,
		ProjectID:                 input.ProjectID,
		ConnectionStringEncrypted: string(encryptedURL),
		SSLMode:                   input.SSLMode,
	}

	err = repo.Create(ctx, environment)

	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	return nil
}
