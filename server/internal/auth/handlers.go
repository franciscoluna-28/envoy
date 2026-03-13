package auth

import (
	"net/http"
	response "newserver/internal/shared"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo      Repository
	validator *validator.Validate
}

func NewHandler(repo Repository, v *validator.Validate) *Handler {
	return &Handler{
		repo:      repo,
		validator: v,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration request"
// @Success 201 {object} User
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	user, err := RegisterUser(r.Context(), h.repo, req.Email, req.Password)
	if err != nil {
		response.WriteJSON(w, http.StatusConflict, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusCreated, user)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} User
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := response.ParseAndValidate(r, h.validator, &req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	user, err := LoginUser(r.Context(), h.repo, req.Email, req.Password)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: err.Error()})
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}
