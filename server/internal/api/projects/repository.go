package projects

import "context"

type Repository interface {
	CreateProject(ctx context.Context, project *Project) error
	GetProjectByID(ctx context.Context, id string) (*Project, error)
	GetProjectsByUserID(ctx context.Context, userID string) ([]*Project, error)
	UpdateProject(ctx context.Context, project *Project) error
	DeleteProject(ctx context.Context, id string) error
}
