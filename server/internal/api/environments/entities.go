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
