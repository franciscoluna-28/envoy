package infra

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/franciscoluna/envoy/server/internal/shared"
)

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

// Register godoc
// @Summary      Register user
// @Description  Creates a new user within the system
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      application.AuthRequest  true  "Credentials"
// @Success      201   {object}  shared.APIResponse{data=application.UserResponse}
// @Failure      400   {object}  shared.APIResponse
// @Failure      409   {object}  shared.APIResponse
// @Router       /auth/register [post]
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request, body application.AuthRequest) error {
	user, err := h.useCase.Execute(r.Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return err
	}

	isProd := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

	SetAuthCookie(w, token, isProd)

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

// Login godoc
// @Summary      Login user
// @Description  Grants access to an existing user within the system
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      application.AuthRequest  true  "Credentials"
// @Success      200   {object}  shared.APIResponse{data=application.UserResponse}
// @Failure      401   {object}  shared.APIResponse
// @Router       /auth/login [post]
func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request, body application.AuthRequest) error {
	user, err := h.useCase.Execute(r.Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return err
	}

	isProd := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	SetAuthCookie(w, token, isProd)

	shared.Success(w, http.StatusOK, application.NewUserResponse(user), "Login successful")
	return nil
}

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value("user_id")

	shared.Success(w, http.StatusOK, map[string]any{
		"id": userID,
	}, "User profile retrieved")

	return nil
}
