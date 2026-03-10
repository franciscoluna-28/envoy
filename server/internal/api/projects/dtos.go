package projects

type ProjectDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type UpdateProjectRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}
