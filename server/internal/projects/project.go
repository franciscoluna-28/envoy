package projects

import "time"

type Project struct {
	ID        string    `json:"id"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name string `validate:"required,min=3,max=100"`
}

type UpdateProjectRequest struct {
	ID   string `validate:"required,uuid4"`
	Name string `validate:"required,min=3,max=100"`
}
