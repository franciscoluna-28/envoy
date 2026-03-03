package infra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/franciscoluna/envoy/server/internal/shared"
)

type RegisterInput struct {
	Body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
}

type RegisterOutput struct {
	Body struct {
		Status  int                      `json:"status"`
		Message string                   `json:"message"`
		Data    application.UserResponse `json:"data,omitempty"`
		Success bool                     `json:"success"`
	}
	SetCookie http.Cookie `header:"Set-Cookie,omitempty"`
}

type LoginInput struct {
	Body struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
}

type LoginOutput struct {
	Body struct {
		Status  int                      `json:"status"`
		Message string                   `json:"message"`
		Data    application.UserResponse `json:"data,omitempty"`
		Success bool                     `json:"success"`
	}
	SetCookie http.Cookie `header:"Set-Cookie,omitempty"`
}

type MeOutput struct {
	Body struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
		Success bool        `json:"success"`
	}
}

type ErrorOutput struct {
	Body struct {
		Status  int                      `json:"status"`
		Message string                   `json:"message"`
		Errors  []map[string]interface{} `json:"errors,omitempty"`
		Success bool                     `json:"success"`
	}
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

func (h *RegisterHandler) Handle(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	user, err := h.useCase.Execute(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		// Handle AppError properly for Huma
		if appErr, ok := err.(*shared.AppError); ok {
			return nil, huma.NewError(appErr.Status, appErr.Msg)
		}
		return nil, err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Create the cookie
	cookie := http.Cookie{
		Name:     auth.AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false during development, true in production
		SameSite: http.SameSiteLaxMode,
	}

	return &RegisterOutput{
		Body: struct {
			Status  int                      `json:"status"`
			Message string                   `json:"message"`
			Data    application.UserResponse `json:"data,omitempty"`
			Success bool                     `json:"success"`
		}{
			Status:  http.StatusCreated,
			Message: "User registered successfully",
			Data:    application.NewUserResponse(user),
			Success: true,
		},
		SetCookie: cookie,
	}, nil
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

func (h *LoginHandler) Handle(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := h.useCase.Execute(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		// Handle AppError properly for Huma
		if appErr, ok := err.(*shared.AppError); ok {
			return nil, huma.NewError(appErr.Status, appErr.Msg)
		}
		return nil, err
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Create the cookie
	cookie := http.Cookie{
		Name:     auth.AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false during development, true in production
		SameSite: http.SameSiteLaxMode,
	}

	return &LoginOutput{
		Body: struct {
			Status  int                      `json:"status"`
			Message string                   `json:"message"`
			Data    application.UserResponse `json:"data,omitempty"`
			Success bool                     `json:"success"`
		}{
			Status:  http.StatusOK,
			Message: "Login successful",
			Data:    application.NewUserResponse(user),
			Success: true,
		},
		SetCookie: cookie,
	}, nil
}

type UserHandler struct {
	authService auth.TokenProvider
}

func NewUserHandler(authService auth.TokenProvider) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

func (h *UserHandler) Handle(ctx context.Context, input *struct{}) (*MeOutput, error) {
	// Try to get the auth cookie from the request
	// In Huma v2, we need to access the context to get cookies
	// For now, let's check if user_id is in context (set by middleware)
	userID := ctx.Value("user_id")
	if userID == nil {
		// If no middleware, try to validate the token manually
		// Note: Huma v2 doesn't provide direct cookie access in the standard context
		// We'll need to handle this at the middleware level
		return nil, fmt.Errorf("unauthorized: no user context found")
	}

	return &MeOutput{
		Body: struct {
			Status  int         `json:"status"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
			Success bool        `json:"success"`
		}{
			Status:  http.StatusOK,
			Message: "User profile retrieved",
			Data:    map[string]any{"id": userID},
			Success: true,
		},
	}, nil
}
