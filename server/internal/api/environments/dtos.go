package environments

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
