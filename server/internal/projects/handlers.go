package projects

import (
	"net/http"
	"time"

	response "newserver/internal/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	repo      Repository
	validator *validator.Validate
}

func NewHandler(repo Repository, validator *validator.Validate) *Handler {
	return &Handler{
		repo:      repo,
		validator: validator,
	}
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProjectRequest true "Project creation request"
// @Success 201 {object} Project
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/projects [post]
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	userID := r.Context().Value("user_id").(string)
	now := time.Now()

	project := Project{
		ID:        uuid.New().String(),
		CreatedBy: userID,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := CreateProject(r.Context(), h.repo, project)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusCreated, project)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update an existing project for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body UpdateProjectRequest true "Project update request"
// @Success 200 {object} Project
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id} [put]
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	req := UpdateProjectRequest{
		ID: chi.URLParam(r, "id"),
	}

	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	userID := r.Context().Value("user_id").(string)

	uProject := Project{
		CreatedBy: userID,
		Name:      req.Name,
		ID:        req.ID,
		UpdatedAt: time.Now(),
	}

	err := UpdateProject(r.Context(), h.repo, uProject, userID)

	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, uProject)
}

// GetProject godoc
// @Summary Get a project
// @Description Get a specific project by ID for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} Project
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id} [get]
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	userID := r.Context().Value("user_id").(string)

	project, err := GetProjectByID(r.Context(), h.repo, projectID, userID)
	if err != nil {
		if err == ErrProjectNotFound {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse{Message: "Project not found"})
			return
		}
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, project)
}

// GetAllProjects godoc
// @Summary Get all projects
// @Description Get all projects for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} Project
// @Failure 500 {object} response.ErrorResponse
// @Router /projects [get]
func (h *Handler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	projects, err := GetAllProjectsByUserID(r.Context(), h.repo, userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, projects)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project by ID for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 204
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id} [delete]
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	userID := r.Context().Value("user_id").(string)

	err := DeleteProject(r.Context(), h.repo, projectID, userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusNoContent, nil)
}
