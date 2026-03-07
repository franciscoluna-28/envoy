package environments

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/shared"
)

// TestConnection godoc
// @Summary      Test database connection
// @Description  Tests if a database connection is valid and has migration permissions
// @Tags         Database
// @Accept       json
// @Produce      json
// @Param        request body DatabaseConnection true "Connection details"
// @Success      200  {object}  shared.APIResponse
// @Failure      400  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/database/test-connection [post]
func HandleTestConnection(w http.ResponseWriter, r *http.Request) error {
	var payload DatabaseConnection

	if err := shared.Decode(r, &payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Invalid JSON body", nil)
		return nil
	}

	if err := shared.Validate(&payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Validation failed", shared.MapValidationErrors(err))
		return nil
	}

	if err := ValidateConnectionAsMigrator(r.Context(), payload); err != nil {
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

	shared.Success(w, http.StatusOK, nil, "Connection successfully tested")
	return nil
}

// CreateEnvironment godoc
// @Summary      Create environment
// @Description  Creates a new database environment for a project
// @Tags         Environments
// @Accept       json
// @Produce      json
// @Param        request body CreateEnvironmentRequest true "Environment details"
// @Success      201  {object}  shared.APIResponse
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Security     CookieAuth
// @Router       /api/v1/environments [post]
func HandleCreateEnvironment(w http.ResponseWriter, r *http.Request) error {
	var payload CreateEnvironmentRequest

	if err := shared.Decode(r, &payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Invalid JSON body", nil)
		return nil
	}

	if err := shared.Validate(&payload); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Validation failed", shared.MapValidationErrors(err))
		return nil
	}

	// TODO: Get master key from context
	masterKey := []byte("master-key")

	if err := CreateProjectEnvironment(r.Context(), payload, masterKey, &EnvRepo{}); err != nil {
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

	shared.Success(w, http.StatusCreated, nil, "Environment created successfully")
	return nil
}
