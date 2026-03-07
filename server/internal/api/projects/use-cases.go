package projects

import "context"

func CreateProject(ctx context.Context, payload CreateProjectRequest, repo *ProjectRepo, userID string) error {
	if err := repo.Create(ctx, payload.Name, userID); err != nil {
		return err
	}
	return nil
}

// You think I'll blunder the security? No. I'm better than that.
func DeleteProject(ctx context.Context, id string, repo *ProjectRepo, userID string) error {
	if err := repo.Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}

func GetProject(ctx context.Context, id string, repo *ProjectRepo, userID string) (*ProjectDTO, error) {
	project, err := repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return &ProjectDTO{
		ID:   project.ID,
		Name: project.Name,
	}, nil
}

func UpdateProject(ctx context.Context, id string, name string, repo *ProjectRepo, userID string) error {
	if err := repo.Update(ctx, id, name, userID); err != nil {
		return err
	}
	return nil
}

func ListProjectsByUser(ctx context.Context, repo *ProjectRepo, userID string) ([]ProjectDTO, error) {
	projects, err := repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var projectDTOs []ProjectDTO
	for _, project := range projects {
		projectDTOs = append(projectDTOs, ProjectDTO{
			ID:   project.ID,
			Name: project.Name,
		})
	}
	return projectDTOs, nil
}
