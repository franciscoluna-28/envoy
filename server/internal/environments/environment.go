package environments

import "time"

type TypeofEnvironment string

const (
	Development TypeofEnvironment = "development"
	Staging     TypeofEnvironment = "staging"
	Production  TypeofEnvironment = "production"
)

type Environment struct {
	ID                        string            `json:"id" db:"id"`
	Name                      string            `json:"name" db:"name"`
	ProjectID                 string            `json:"project_id" db:"project_id"`
	Type                      TypeofEnvironment `json:"type" db:"type"`
	ConnectionStringEncrypted string            `json:"-" db:"connection_string_encrypted"`
	ConnectionError           string            `json:"connection_error" db:"connection_error"`
	CreatedAt                 time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time         `json:"updated_at" db:"updated_at"`
}

type EnvironmentDbUser struct {
	ID                        string    `json:"id" db:"id"`
	EnvironmentID             string    `json:"environment_id" db:"environment_id"`
	Name                      string    `json:"name" db:"name"`
	ConnectionStringEncrypted string    `json:"-" db:"connection_string_encrypted"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

type DatabaseConnection struct {
	ConnectionString string `json:"connection_string" validate:"required"`
}

type CreateEnvironmentRequest struct {
	Name          string            `json:"name" validate:"required"`
	Type          TypeofEnvironment `json:"type" validate:"required,oneof=development staging production"`
	ProjectID     string            `json:"project_id" validate:"required"`
	ConnectionUrl string            `json:"connection_url" validate:"required"`
}

type CertificatesConfig struct {
	CACert     string `json:"ca_cert,omitempty"`
	ClientCert string `json:"client_cert,omitempty"`
	ClientKey  string `json:"client_key,omitempty"`
}

type SchemaColumn struct {
	TableName  string `json:"table_name"`
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	IsNullable string `json:"is_nullable"`
}

type EnvironmentMigration struct {
	ID            string     `json:"id" db:"id"`
	EnvironmentID string     `json:"environment_id" db:"environment_id"`
	Name          string     `json:"name" db:"name"`
	Description   string     `json:"description" db:"description"`
	SQLContent    string     `json:"sql_content" db:"sql_content"`
	Status        string     `json:"status" db:"status"`
	ExecutedAt    *time.Time `json:"executed_at" db:"executed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	Duration      int64      `json:"duration" db:"duration"` // Duration in milliseconds to match INTEGER type
	ExecutedBy    string     `json:"executed_by" db:"executed_by"`
	ErrorMessage  string     `json:"error_message" db:"error_message"`
	ClientId      string     `json:"client_id" db:"client_id"`
	Checksum      string     `json:"-" db:"checksum"`
}

// Migration status constants
const (
	MigrationStatusPending   string = "pending"
	MigrationStatusRunning   string = "running"
	MigrationStatusCompleted string = "completed"
	MigrationStatusFailed    string = "failed"
)

type CreateEnvironmentMigrationRequest struct {
	EnvironmentID string `json:"environment_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SQLContent    string `json:"sql_content" validate:"required"`
	ClientId      string `json:"client_id" validate:"required"`
}

type PreviewMigrationRequest struct {
	SQLContent string `json:"sql_content" validate:"required"`
}

type UpdateEnvironmentRequest struct {
	Name          string            `json:"name" validate:"required"`
	Type          TypeofEnvironment `json:"type" validate:"required,oneof=development staging production"`
	ConnectionUrl string            `json:"connection_url" validate:"required"`
}
