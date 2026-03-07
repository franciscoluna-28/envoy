package projects

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(dbPath string) (*ProjectRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return &ProjectRepo{db: db}, nil
}

func (r *ProjectRepo) Create(ctx context.Context, name, createdBy string) error {
	query := `INSERT INTO projects (id, name, created_by, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), name, createdBy, time.Now())
	return err
}

// Multi tenant systems taught me to be paranoid
func (r *ProjectRepo) GetByID(ctx context.Context, id string, userID string) (*Project, error) {
	query := `SELECT id, name, created_by, created_at FROM projects WHERE id = ? AND created_by = ?`
	var project Project
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(&project.ID, &project.Name, &project.CreatedBy, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) ListByUser(ctx context.Context, userID string) ([]Project, error) {
	query := `SELECT id, name, created_by, created_at FROM projects WHERE created_by = ?`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Name, &project.CreatedBy, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM projects WHERE id = ? AND created_by = ?`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}

func (r *ProjectRepo) Update(ctx context.Context, id, name string, userID string) error {
	query := `UPDATE projects SET name = ?, updated_at = ? WHERE id = ? AND created_by = ?`
	_, err := r.db.ExecContext(ctx, query, name, time.Now(), id, userID)
	return err
}
