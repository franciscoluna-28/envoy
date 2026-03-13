package projects

import "context"

func CreateProject(ctx context.Context, repo Repository, project Project) error {
	return repo.Create(ctx, project)
}

func GetProjectByID(ctx context.Context, repo Repository, id string, userID string) (*Project, error) {
	return repo.GetByID(ctx, id, userID)
}

func GetAllProjectsByUserID(ctx context.Context, repo Repository, userID string) ([]Project, error) {
	return repo.GetAllByUserID(ctx, userID)
}

func UpdateProject(ctx context.Context, repo Repository, project Project, userID string) error {
	return repo.Update(ctx, project, userID)
}

func DeleteProject(ctx context.Context, repo Repository, id string, userID string) error {
	return repo.Delete(ctx, id, userID)
}
