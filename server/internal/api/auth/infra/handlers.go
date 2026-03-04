package infra

import (
	"net/http"
	"time"

	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
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

type RegisterHandler struct {
	useCase     *application.RegisterWithEmailUseCase
	authService auth.TokenProvider
}

func NewRegisterHandler(uc *application.RegisterWithEmailUseCase, auth auth.TokenProvider) *RegisterHandler {
	return &RegisterHandler{
		useCase:     uc,
		authService: auth,
	}
}

// Handle godoc
// @Summary      Register user
// @Description  Creates a new user within the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterInput true "Request body"
// @Success      201  {object}  shared.APIResponse{data=application.UserResponse}
// @Failure      400  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/auth/register [post]
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var input RegisterInput
	if err := shared.Decode(r, &input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, nil)
		return nil
	}

	if err := shared.Validate(&input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, nil)
		return nil
	}

	user, err := h.useCase.Execute(r.Context(), input.Email, input.Password)
	if err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		shared.InternalError(w)
		return nil
	}

	cookie := http.Cookie{
		Name:     auth.AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false during development, true in production
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	shared.Success(w, http.StatusCreated, application.NewUserResponse(user), "User registered successfully")
	return nil
}

type LoginHandler struct {
	useCase     *application.LoginWithEmailUseCase
	authService auth.TokenProvider
}

func NewLoginHandler(uc *application.LoginWithEmailUseCase, auth auth.TokenProvider) *LoginHandler {
	return &LoginHandler{
		useCase:     uc,
		authService: auth,
	}
}

// Handle godoc
// @Summary      Login user
// @Description  Grants access to an existing user within the system
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginInput true "Request body"
// @Success      200  {object}  shared.APIResponse{data=application.UserResponse}
// @Failure      400  {object}  shared.APIResponse
// @Failure      401  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/auth/login [post]
func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var input LoginInput
	if err := shared.Decode(r, &input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, nil)
		return nil
	}

	if err := shared.Validate(&input); err != nil {
		shared.BadRequest(w, shared.ErrInvalidInput, nil)
		return nil
	}

	user, err := h.useCase.Execute(r.Context(), input.Email, input.Password)
	if err != nil {
		if appErr, ok := err.(*shared.AppError); ok {
			shared.Send(w, appErr.Status, shared.APIResponse{
				Status:  appErr.Status,
				Message: appErr.Msg,
				Success: false,
			})
			return nil
		}
		shared.InternalError(w)
		return nil
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		shared.InternalError(w)
		return nil
	}

	cookie := http.Cookie{
		Name:     auth.AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false during development, true in production
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	shared.Success(w, http.StatusOK, application.NewUserResponse(user), "Login successful")
	return nil
}

type UserHandler struct {
	authService auth.TokenProvider
}

func NewUserHandler(authService auth.TokenProvider) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

// LogoutHandler handles clearing the authentication cookie on logout.
type LogoutHandler struct {
}

func NewLogoutHandler() *LogoutHandler {
	return &LogoutHandler{}
}

// Handle godoc
// @Summary      Logout user
// @Description  Clears authentication cookie to log the user out
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  shared.APIResponse
// @Failure      500  {object}  shared.APIResponse
// @Router       /api/v1/auth/logout [post]
func (h *LogoutHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	// Overwrite the auth cookie with an expired one
	cookie := http.Cookie{
		Name:     auth.AuthCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // keep same flags as other cookies
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	http.SetCookie(w, &cookie)

	shared.Success(w, http.StatusOK, nil, "Logout successful")
	return nil
}

// Handle godoc
// @Summary      Get current user
// @Description  Get the current authenticated user's profile
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  shared.APIResponse{data=object}
// @Failure      401  {object}  shared.APIResponse
// @Router       /api/v1/me [get]
func (h *UserHandler) Handle(w http.ResponseWriter, r *http.Request) error {
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
