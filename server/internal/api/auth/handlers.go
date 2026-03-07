package auth

import (
	"net/http"
	"time"

	"github.com/franciscoluna/envoy/server/internal/shared"
)

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// HandleRegister godoc
// @Summary      Register user
// @Description  Creates a new user within the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterInput true "Request body"
// @Success      201  {object}  shared.APIResponse{data=UserResponse}
// @Failure      400  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/auth/register [post]
func HandleRegister(w http.ResponseWriter, r *http.Request, repo UserRepository, tp TokenProvider) error {
	var input RegisterInput
	if err := shared.Decode(r, &input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Invalid JSON body", nil)
		return nil
	}
	if err := shared.Validate(&input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Validation failed", shared.MapValidationErrors(err))
		return nil
	}

	user, err := RegisterUser(r.Context(), repo, input.Email, input.Password)
	if err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{Status: appErr.Status, Message: appErr.Msg, Success: false})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	token, err := tp.GenerateToken(user)
	if err != nil {
		shared.InternalError(w)
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	shared.Success(w, http.StatusCreated, userResponse(user), "User registered successfully")
	return nil
}

// HandleLogin godoc
// @Summary      Login user
// @Description  Grants access to an existing user within the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginInput true "Request body"
// @Success      200  {object}  shared.APIResponse{data=UserResponse}
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/auth/login [post]
func HandleLogin(w http.ResponseWriter, r *http.Request, repo UserRepository, tp TokenProvider) error {
	var input LoginInput
	if err := shared.Decode(r, &input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Invalid JSON body", nil)
		return nil
	}
	if err := shared.Validate(&input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, "Validation failed", shared.MapValidationErrors(err))
		return nil
	}

	user, err := LoginUser(r.Context(), repo, input.Email, input.Password)
	if err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{Status: appErr.Status, Message: appErr.Msg, Success: false})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	token, err := tp.GenerateToken(user)
	if err != nil {
		shared.InternalError(w)
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	shared.Success(w, http.StatusOK, userResponse(user), "Login successful")
	return nil
}

// HandleLogout godoc
// @Summary      Logout user
// @Description  Clears authentication cookie to log the user out
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/auth/logout [post]
func HandleLogout(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	shared.Success(w, http.StatusOK, nil, "Logout successful")
	return nil
}

// HandleMe godoc
// @Summary      Get current user
// @Description  Get the current authenticated user's profile
// @Tags         User
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  shared.APIResponse{data=object}
// @Failure      401  {object}  shared.APIResponse
// @Router       /api/v1/me [get]
func HandleMe(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value("user_id")
	if userID == nil {
		shared.Send(w, http.StatusUnauthorized, shared.APIResponse{
			Status:  http.StatusUnauthorized,
			Error:   shared.ErrUnauthorized,
			Message: "unauthorized: no user context found",
			Success: false,
		})
		return nil
	}
	shared.Success(w, http.StatusOK, map[string]any{"id": userID}, "User profile retrieved")
	return nil
}
