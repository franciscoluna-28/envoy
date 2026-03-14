package environments

import "time"

type Environment struct {
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	ProjectID                 string    `json:"project_id"`
	ConnectionStringEncrypted string    `json:"-"`
	SSLMode                   string    `json:"-"`
	CertificatesJSON          []string  `json:"certificates_json"`
	ConnectionStatus          string    `json:"connection_status"`
	ConnectionError           string    `json:"connection_error"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type EnvironmentDbUser struct {
	ID                        string    `json:"id"`
	EnvironmentID             string    `json:"environment_id"`
	Name                      string    `json:"name"`
	ConnectionStringEncrypted string    `json:"-"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
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
	SSLMode          string `json:"ssl_mode"`
	CertificatesJSON string `json:"certificates_json"`
}

type CreateEnvironmentRequest struct {
	Name          string              `json:"name" validate:"required"`
	ProjectID     string              `json:"project_id" validate:"required"`
	ConnectionUrl string              `json:"connection_url" validate:"required"`
	SSLMode       string              `json:"ssl_mode" validate:"oneof=disable allow prefer require verify-ca verify-full"`
	Certificates  *CertificatesConfig `json:"certificates,omitempty"`
}

type CertificatesConfig struct {
	CACert     string `json:"ca_cert,omitempty"`
	ClientCert string `json:"client_cert,omitempty"`
	ClientKey  string `json:"client_key,omitempty"`
}