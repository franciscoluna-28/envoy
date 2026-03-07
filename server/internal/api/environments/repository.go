package environments

import "context"

type Repository interface {
	CreateEnvironment(ctx context.Context, env *Environment) error
	GetEnvironmentByID(ctx context.Context, id string) (*Environment, error)
	GetEnvironmentsByProjectID(ctx context.Context, projectID string) ([]*Environment, error)
	UpdateEnvironment(ctx context.Context, env *Environment) error
	DeleteEnvironment(ctx context.Context, id string) error
}
