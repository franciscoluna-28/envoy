package handlers

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/api/dto"
	"github.com/franciscoluna/envoy/server/internal/api/infra"
	"github.com/franciscoluna/envoy/server/internal/api/usecases"
	"github.com/franciscoluna/envoy/server/internal/shared"
)

type RegisterHandler struct {
	useCase     *usecases.RegisterWithEmailUseCase
	authService *infra.AuthService
}

func NewRegisterHandler(uc *usecases.RegisterWithEmailUseCase, auth *infra.AuthService) *RegisterHandler {
	return &RegisterHandler{
		useCase:     uc,
		authService: auth,
	}
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var body dto.AuthRequest

	if err := infra.Decode(r, &body); err != nil {
		return err
	}

	if err := infra.Validate(body); err != nil {
		shared.BadRequest(w, "invalid_input", infra.MapValidationErrors(err))
		return nil
	}

	user, err := h.useCase.Execute(r.Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return err
	}

	isProd := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	infra.SetAuthCookie(w, token, isProd)

	shared.Success(w, http.StatusCreated, dto.NewUserResponse(user), "user registered successfully")
	return nil
}

type LoginHandler struct {
	useCase     *usecases.LoginWithEmailUseCase
	authService *infra.AuthService
}

func NewLoginHandler(uc *usecases.LoginWithEmailUseCase, auth *infra.AuthService) *LoginHandler {
	return &LoginHandler{
		useCase:     uc,
		authService: auth,
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var body dto.AuthRequest

	if err := infra.Decode(r, &body); err != nil {
		return err
	}

	if err := infra.Validate(body); err != nil {
		shared.BadRequest(w, "invalid_credentials_format", infra.MapValidationErrors(err))
		return nil
	}
	user, err := h.useCase.Execute(r.Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return err
	}

	isProd := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	infra.SetAuthCookie(w, token, isProd)

	shared.Success(w, http.StatusOK, dto.NewUserResponse(user), "Login successful")
	return nil
}
