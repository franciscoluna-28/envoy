package environments

import (
	"fmt"
	"net/http"
	"newserver/internal/auth"
	response "newserver/internal/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo        Repository
	validator   *validator.Validate
	masterKey   []byte
	checksumKey []byte
}

func NewHandler(repo Repository, v *validator.Validate, masterKey, checksumKey []byte) *Handler {
	return &Handler{
		repo:        repo,
		validator:   v,
		masterKey:   masterKey,
		checksumKey: checksumKey,
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

	var req CreateEnvironmentRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	req.ProjectID = projectID

	err := CreateProjectEnvironment(r.Context(), req, h.masterKey, h.repo)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

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

	envs, err := GetAllEnvironmentsByProjectID(r.Context(), projectID, h.repo)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get environments"})
		return
	}

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

	schema, err := GetEnvironmentSchema(r.Context(), envID, h.repo, h.masterKey)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get environment schema"})
		return
	}

	response.WriteJSON(w, http.StatusOK, schema)
}

// PreviewEnvironmentSchemaChanges godoc
// @Summary Preview environment schema changes
// @Description Preview database schema changes by running migration in a transaction
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Param request body PreviewMigrationRequest true "Migration preview request"
// @Success 200 {array} SchemaColumn
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/migrations/preview [post]
func (h *Handler) PreviewEnvironmentSchemaChanges(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	var req PreviewMigrationRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	schema, err := PreviewEnvironmentSchemaChanges(r.Context(), envID, h.repo, req.SQLContent, h.masterKey)
	if err != nil {
		// Return detailed error information instead of generic 500
		response.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"errors":  []string{err.Error()},
			"message": "SQL validation failed",
		})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    schema,
		"message": "Schema preview completed successfully",
	})
}

// RunDatabaseMigration godoc
// @Summary Run database migration
// @Description Execute a database migration in an environment with full tracking
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Param request body CreateEnvironmentMigrationRequest true "Migration execution request"
// @Success 201 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/migrations [post]
func (h *Handler) RunDatabaseMigration(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	fmt.Printf("[DEBUG] RunDatabaseMigration handler called with envID: %s\n", envID)

	var req CreateEnvironmentMigrationRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		fmt.Printf("[DEBUG] Failed to parse request: %v\n", err)
		response.WriteValidationError(w, err)
		return
	}

	fmt.Printf("[DEBUG] Request parsed successfully: %+v\n", req)

	val := r.Context().Value(auth.UserIDKey)
	fmt.Printf("[DEBUG] Context value for user_id: %v (type: %T)\n", val, val)

	userID, ok := val.(string)
	if !ok {
		fmt.Printf("[DEBUG] Failed to extract user_id from context: ok=%v, val=%v\n", ok, val)
		response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Unauthorized: user_id missing"})
		return
	}

	fmt.Printf("[DEBUG] Successfully extracted userID: %s\n", userID)

	req.EnvironmentID = envID

	err := RunDatabaseMigrationInEnvironment(r.Context(), envID, h.repo, req, userID, h.masterKey, h.checksumKey)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.ErrorResponse{Message: "Migration executed successfully"})
}

// GetEnvironmentMigrations godoc
// @Summary Get environment migrations
// @Description Get all migrations for a specific environment
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Success 200 {array} EnvironmentMigration
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/migrations [get]
func (h *Handler) GetEnvironmentMigrations(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	migrations, err := h.repo.GetEnvironmentMigrations(r.Context(), envID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get migrations"})
		return
	}

	response.WriteJSON(w, http.StatusOK, migrations)
}

func (h *Handler) VerifyDatabasePermissions(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	err := VerifyDatabasePermissions(r.Context(), envID, h.repo, h.masterKey)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ErrorResponse{Message: "Database permissions verified successfully"})
}

// GetEnvironmentMigrationByID godoc
// @Summary Get environment migration by ID
// @Description Get a specific migration by ID
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Migration ID"
// @Success 200 {object} EnvironmentMigration
// @Failure 500 {object} response.ErrorResponse
// @Router /migrations/{id} [get]
func (h *Handler) GetEnvironmentMigrationByID(w http.ResponseWriter, r *http.Request) {
	migrationID := chi.URLParam(r, "id")

	migration, err := GetEnvironmentMigrationByID(r.Context(), migrationID, h.repo)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to get migration"})
		return
	}

	response.WriteJSON(w, http.StatusOK, migration)
}

// UpdateEnvironment godoc
// @Summary Update environment
// @Description Update an existing environment
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Param request body UpdateEnvironmentRequest true "Environment update request"
// @Success 200 {object} Environment
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id} [put]
func (h *Handler) UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	var req UpdateEnvironmentRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	env, err := UpdateEnvironment(r.Context(), envID, req, h.masterKey, h.repo)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, env)
}

// DeleteEnvironment godoc
// @Summary Delete environment
// @Description Delete an environment
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Success 200 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id} [delete]
func (h *Handler) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	err := h.repo.DeleteEnvironment(r.Context(), envID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to delete environment"})
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ErrorResponse{Message: "Environment deleted successfully"})
}

// ValidateEnvironmentConnection godoc
// @Summary Validate environment connection
// @Description Validate database connection for an environment
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Success 200 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/validate [post]
func (h *Handler) ValidateEnvironmentConnection(w http.ResponseWriter, r *http.Request) {
	if err := ValidateEnvironmentConnection(r.Context(), chi.URLParam(r, "id"), h.repo, h.masterKey); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ErrorResponse{Message: "Environment connection validated successfully"})
}

// TestPermissionsWithPreview godoc
// @Summary Test permissions with preview schema
// @Description Test database user permissions against schema changes in a transaction
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Param request body TestPermissionsRequest true "Permission test request"
// @Success 200 {array} TablePermission
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/test-permissions-preview [post]
func (h *Handler) TestPermissionsWithPreview(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	var req TestPermissionsRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	permissions, err := AuditPermissionsBeforeMigration(r.Context(), envID, h.repo, req.SQLContent, req.DatabaseUser, h.masterKey)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    permissions,
		"message": "Permission test completed successfully",
	})
}

// TestPermissionsWithCurrentSchema godoc
// @Summary Test permissions with current schema
// @Description Test database user permissions against current schema
// @Tags environments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Environment ID"
// @Param request body TestPermissionsWithCurrentSchemaRequest true "Permission test request"
// @Success 200 {array} TablePermission
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /environments/{id}/test-permissions-current [post]
func (h *Handler) TestPermissionsWithCurrentSchema(w http.ResponseWriter, r *http.Request) {
	envID := chi.URLParam(r, "id")

	var req TestPermissionsWithCurrentSchemaRequest
	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	permissions, err := AuditPermissionsWithCurrentSchema(r.Context(), envID, h.repo, req.DatabaseUser, h.masterKey)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    permissions,
		"message": "Permission test completed successfully",
	})
}
