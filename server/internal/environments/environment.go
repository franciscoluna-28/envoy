package environments

import "time"

type Environment struct {
	ID                        string    `json:"id" db:"id"`
	Name                      string    `json:"name" db:"name"`
	ProjectID                 string    `json:"project_id" db:"project_id"`
	ConnectionStringEncrypted string    `json:"-" db:"connection_string_encrypted"`
	ConnectionError           string    `json:"connection_error" db:"connection_error"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

type EnvironmentDbUser struct {
	ID                        string    `json:"id" db:"id"`
	EnvironmentID             string    `json:"environment_id" db:"environment_id"`
	Name                      string    `json:"name" db:"name"`
	ConnectionStringEncrypted string    `json:"-" db:"connection_string_encrypted"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

type EnvironmentMigration struct {
	ID            string    `json:"id"`
	EnvironmentID string    `json:"environment_id"`
	Name          string    `json:"name"`
	SQLContent    string    `json:"sql_content"`
	Status        string    `json:"status"`
	ExecutedAt    time.Time `json:"executed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type DatabaseConnection struct {
	ConnectionString string `json:"connection_string" validate:"required"`
}

type CreateEnvironmentRequest struct {
	Name          string `json:"name" validate:"required"`
	ProjectID     string `json:"project_id" validate:"required"`
	ConnectionUrl string `json:"connection_url" validate:"required"`
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
