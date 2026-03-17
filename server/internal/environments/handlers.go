package environments

import (
	"fmt"
	"net/http"
	response "newserver/internal/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo      Repository
	validator *validator.Validate
	key       []byte
}

func NewHandler(repo Repository, v *validator.Validate, key []byte) *Handler {
	return &Handler{
		repo:      repo,
		validator: v,
		key:       key,
	}
}

// CreateEnvironment godoc
// @Summary Create a new project environment
// @Description Create a new project environment for the authenticated user
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body CreateEnvironmentRequest true "Environment creation request"
// @Success 201 {object} Environment
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id}/environments [post]
func (h *Handler) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	fmt.Printf("CreateEnvironment: Handler called for project ID: %s\n", projectID)

	var req CreateEnvironmentRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		fmt.Printf("CreateEnvironment: Validation error - %v\n", err)
		response.WriteValidationError(w, err)
		return
	}

	req.ProjectID = projectID
	fmt.Printf("CreateEnvironment: Creating environment with name: %s\n", req.Name)

	err := CreateProjectEnvironment(r.Context(), req, h.key, h.repo)
	if err != nil {
		fmt.Printf("CreateEnvironment: Error creating environment - %v\n", err)
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	fmt.Printf("CreateEnvironment: Successfully created environment\n")
	response.WriteJSON(w, http.StatusCreated, response.ErrorResponse{Message: "Environment created successfully"})
}

// GetAllEnvironmentsByProjectID godoc
// @Summary Get all environments for project
// @Description Get all environments for a specific project
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {array} Environment
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id}/environments [get]
func (h *Handler) GetAllEnvironmentsByProjectID(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	fmt.Printf("GetAllEnvironmentsByProjectID: Handler called for project ID: %s\n", projectID)

	envs, err := GetAllEnvironmentsByProjectID(r.Context(), projectID, h.repo)
	if err != nil {
		fmt.Printf("GetAllEnvironmentsByProjectID: Error getting environments - %v\n", err)
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get environments"})
		return
	}

	fmt.Printf("GetAllEnvironmentsByProjectID: Successfully retrieved %d environments\n", len(envs))
	response.WriteJSON(w, http.StatusOK, envs)
}

// GetEnvironmentByID godoc
// @Summary Get environment by ID
// @Description Get a specific environment by ID
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Success 200 {object} Environment
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id} [get]
func (h *Handler) GetEnvironmentByID(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	env, err := GetEnvironmentByID(r.Context(), envID, h.repo)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get environment"})
		return
	}

	response.WriteJSON(w, http.StatusOK, env)
}

// GetEnvironmentSchema godoc
// @Summary Get environment schema
// @Description Get database schema for a specific environment
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/schema [get]
func (h *Handler) GetEnvironmentSchema(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")
	fmt.Printf("GetEnvironmentSchema: Handler called for environment ID: %s\n", envID)

	schema, err := GetEnvironmentSchema(r.Context(), envID, h.repo, h.key)
	if err != nil {
		fmt.Printf("GetEnvironmentSchema: Error getting schema - %v\n", err)
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get environment schema"})
		return
	}

	fmt.Printf("GetEnvironmentSchema: Successfully retrieved schema for environment %s\n", envID)
	response.WriteJSON(w, http.StatusOK, schema)
}
