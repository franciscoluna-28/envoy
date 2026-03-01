package handlers

import (
	"net/http"
	"time"

	"github.com/franciscoluna/envoy/server/internal/api/domain"
	"github.com/franciscoluna/envoy/server/internal/api/infra"
	"github.com/franciscoluna/envoy/server/internal/api/usecases"
)

func setAuthCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(168 * time.Hour), // 7 days
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

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
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := infra.Decode(r, &body); err != nil {
		return err
	}

	if err := infra.Validate(body); err != nil {
		infra.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": infra.MapValidationErrors(err),
		})
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

	setAuthCookie(w, token, isProd)

	infra.JSON(w, http.StatusCreated, domain.NewUserResponse(user))
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
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := infra.Decode(r, &body); err != nil {
		return err
	}
	if err := infra.Validate(body); err != nil {
		infra.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": infra.MapValidationErrors(err),
		})
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
	setAuthCookie(w, token, isProd)

	infra.JSON(w, http.StatusOK, domain.NewUserResponse(user))
	return nil
}
