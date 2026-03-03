package infra

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
)

func RegisterRoutes(api huma.API, authService auth.TokenProvider, regUC *application.RegisterWithEmailUseCase, loginUC *application.LoginWithEmailUseCase) {
	regHandler := NewRegisterHandler(regUC, authService)
	loginHandler := NewLoginHandler(loginUC, authService)
	userHandler := NewUserHandler(authService)

	huma.Register(api, huma.Operation{
		OperationID: "post-register",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/register",
		Summary:     "Register user",
		Description: "Creates a new user within the system",
		Tags:        []string{"Auth"},
	}, regHandler.Handle)

	huma.Register(api, huma.Operation{
		OperationID: "post-login",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/login",
		Summary:     "Login user",
		Description: "Grants access to an existing user within the system",
		Tags:        []string{"Auth"},
	}, loginHandler.Handle)

	huma.Register(api, huma.Operation{
		OperationID: "get-me",
		Method:      http.MethodGet,
		Path:        "/api/v1/me",
		Summary:     "Get current user",
		Description: "Get the current authenticated user's profile",
		Tags:        []string{"User"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, userHandler.Handle)
}
