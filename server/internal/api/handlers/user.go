package handlers

import (
	"net/http"
	infra "server/internal/api/infra"
	web "server/internal/api/infra"
	"server/internal/api/usecases"
	"time"
)

type RegisterHandler struct {
	useCase     *usecases.RegisterWithEmailUseCase
	authService *infra.AuthService
}

func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) error {
	var body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	isProd := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" // TODO: Improve this implementation

	if err := web.Decode(r, &body); err != nil {
		return err
	}

	user, err := h.useCase.Execute(r.Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	infra.JSON(w, http.StatusCreated, user)
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

	user, err := h.useCase.Execute(r.Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return err
	}

	isProd := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(168 * time.Hour),
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	infra.JSON(w, http.StatusOK, user)
	return nil
}
