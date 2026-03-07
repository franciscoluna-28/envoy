package projects

type ProjectDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateProjectRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}
