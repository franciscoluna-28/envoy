package domain

import "time"

type UserID string
type ProjectID string
type EnvironmentID string
type EnvironmentMigrationID string
type EnvironmentDbUserID string
type MigrationStatus string

const (
	StatusPending  MigrationStatus = "pending"
	StatusExecuted MigrationStatus = "executed"
	StatusFailed   MigrationStatus = "failed"
)

type User struct {
	ID           UserID    `json:"id"            db:"id"`
	Email        string    `json:"email"         db:"email"`
	PasswordHash string    `json:"-"             db:"password_hash"`
	CreatedAt    time.Time `json:"created_at"    db:"created_at"`
}

type Project struct {
	ID        ProjectID `json:"id"`
	CreatedBy UserID    `json:"created_by"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Environment struct {
	ID                        EnvironmentID `json:"id"`
	Name                      string        `json:"name"`
	ProjectID                 ProjectID     `json:"project_id"`
	ConnectionStringEncrypted string        `json:"-"`
	SSLMode                   string        `json:"-"`
	CertificatesJSON          []string      `json:"certificates_json"`
	ConnectionStatus          string        `json:"connection_status"`
	ConnectionError           string        `json:"connection_error"`
	CreatedAt                 time.Time     `json:"created_at"`
	UpdatedAt                 time.Time     `json:"updated_at"`
}

type EnvironmentDbUser struct {
	ID                        EnvironmentDbUserID `json:"id"`
	EnvironmentID             EnvironmentID       `json:"environment_id"`
	Name                      string              `json:"name"`
	ConnectionStringEncrypted string              `json:"-"`
	CreatedAt                 time.Time           `json:"created_at"`
	UpdatedAt                 time.Time           `json:"updated_at"`
}

type EnvironmentMigration struct {
	ID            EnvironmentMigrationID `json:"id"`
	EnvironmentID EnvironmentID          `json:"environment_id"`
	Name          string                 `json:"name"`
	SQLContent    string                 `json:"sql_content"`
	Status        MigrationStatus        `json:"status"`
	ExecutedAt    time.Time              `json:"executed_at"`
	CreatedAt     time.Time              `json:"created_at"`
}
