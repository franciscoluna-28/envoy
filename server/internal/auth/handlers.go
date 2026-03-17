package auth

import (
	"net/http"
	response "newserver/internal/shared"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo        Repository
	validator   *validator.Validate
	jwtProvider *JWTProvider
}

func NewHandler(repo Repository, v *validator.Validate, jwtProvider *JWTProvider) *Handler {
	return &Handler{
		repo:        repo,
		validator:   v,
		jwtProvider: jwtProvider,
	}
}

// setJWTCookie sets the JWT token as an HTTP cookie
func (h *Handler) setJWTCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24 hours in seconds
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
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

	// Generate JWT token and set cookie
	token, err := h.jwtProvider.GenerateToken(string(user.ID))
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to generate token"})
		return
	}
	h.setJWTCookie(w, token)

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

	// Generate JWT token and set cookie
	token, err := h.jwtProvider.GenerateToken(string(user.ID))
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.ErrorResponse{Message: "Failed to generate token"})
		return
	}
	h.setJWTCookie(w, token)

	response.WriteJSON(w, http.StatusOK, user)
}

// GetMe godoc
// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "User not found in context"})
		return
	}

	user, err := h.repo.GetById(r.Context(), userID)

	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "User not found"})
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// Logout godoc
// @Summary User logout
// @Description Logout the current user by clearing the auth cookie
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie by setting it with MaxAge -1
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	response.WriteJSON(w, http.StatusOK, response.ErrorResponse{Message: "Logged out successfully"})
}
