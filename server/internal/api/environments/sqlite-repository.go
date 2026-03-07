package environments

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type EnvRepo struct {
	db *sql.DB
}

func NewEnvRepo(dbPath string) (*EnvRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return &EnvRepo{db: db}, nil
}

func (r *EnvRepo) Create(ctx context.Context, name, projectID, encryptedURL, sslMode string) error {
	query := `INSERT INTO environments (id, name, project_id, connection_string_encrypted, ssl_mode, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		uuid.New().String(),
		name,
		projectID,
		encryptedURL,
		sslMode,
		time.Now(),
		time.Now(),
	)
	return err
}
