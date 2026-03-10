package projects

import (
	"fmt"
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-chi/chi/v5"
)

// CreateProject godoc
// @Summary      Create project
// @Description  Creates a new project for the authenticated user
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        request body CreateProjectRequest true "Project details"
// @Success      201  {object}  shared.APIResponse
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Security     CookieAuth
// @Router       /api/v1/projects [post]
func HandleCreateProject(w http.ResponseWriter, r *http.Request, repo *ProjectRepo) error {
	var payload CreateProjectRequest

	if err := shared.Decode(r, &payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Invalid JSON body", nil)
		return nil
	}

	if err := shared.Validate(&payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Validation failed", shared.MapValidationErrors(err))
		return nil
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		shared.InternalError(w)
		return nil
	}
	if err := CreateProject(r.Context(), payload, repo, userID); err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Error:   appErr.Code,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	shared.Success(w, http.StatusCreated, nil, "Project created successfully")
	return nil
}

// GetProject godoc
// @Summary      Get project by ID
// @Description  Retrieves a specific project belonging to the authenticated user
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id path string true "Project ID"
// @Success      200  {object}  shared.APIResponse{data=ProjectDTO}
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      404  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Security     CookieAuth
// @Router       /api/v1/projects/{id} [get]
func HandleGetProject(w http.ResponseWriter, r *http.Request, repo *ProjectRepo) error {
	id := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		shared.InternalError(w)
		return nil
	}

	project, err := GetProject(r.Context(), id, repo, userID)
	if err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Error:   appErr.Code,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	shared.Success(w, http.StatusOK, project, "Project retrieved successfully")
	return nil
}

// ListProjects godoc
// @Summary      List user projects
// @Description  Retrieves all projects belonging to the authenticated user
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Success      200  {object}  shared.APIResponse{data=[]ProjectDTO}
// @Failure      401  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Security     CookieAuth
// @Router       /api/v1/projects [get]
func HandleListProjects(w http.ResponseWriter, r *http.Request, repo *ProjectRepo) error {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		shared.InternalError(w)
		return nil
	}

	projects, err := ListProjectsByUser(r.Context(), repo, userID)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Database error in ListProjectsByUser: %v, userID: %s\n", err, userID)
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Error:   appErr.Code,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	shared.Success(w, http.StatusOK, projects, "Projects listed successfully")
	return nil
}

// UpdateProject godoc
// @Summary      Update project
// @Description  Updates a project name belonging to the authenticated user
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id path string true "Project ID"
// @Param        request body UpdateProjectRequest true "Updated project details"
// @Success      200  {object}  shared.APIResponse
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      404  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Security     CookieAuth
// @Router       /api/v1/projects/{id} [put]
func HandleUpdateProject(w http.ResponseWriter, r *http.Request, repo *ProjectRepo) error {
	id := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		shared.InternalError(w)
		return nil
	}

	var payload UpdateProjectRequest
	if err := shared.Decode(r, &payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Invalid JSON body", nil)
		return nil
	}

	if err := shared.Validate(&payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Validation failed", shared.MapValidationErrors(err))
		return nil
	}

	if err := UpdateProject(r.Context(), id, payload.Name, repo, userID); err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Error:   appErr.Code,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	shared.Success(w, http.StatusOK, nil, "Project updated successfully")
	return nil
}

// DeleteProject godoc
// @Summary      Delete project
// @Description  Deletes a project belonging to the authenticated user
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id path string true "Project ID"
// @Success      204  {object}  shared.APIResponse
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      404  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Security     CookieAuth
// @Router       /api/v1/projects/{id} [delete]
func HandleDeleteProject(w http.ResponseWriter, r *http.Request, repo *ProjectRepo) error {
	id := chi.URLParam(r, "id")
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		shared.InternalError(w)
		return nil
	}

	if err := DeleteProject(r.Context(), id, repo, userID); err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Error:   appErr.Code,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	shared.Success(w, http.StatusNoContent, nil, "Project deleted successfully")
	return nil
}
