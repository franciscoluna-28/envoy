package infra

import (
	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, authService auth.TokenProvider, regUC *application.RegisterWithEmailUseCase, loginUC *application.LoginWithEmailUseCase) {
	regHandler := NewRegisterHandler(regUC, authService)
	loginHandler := NewLoginHandler(loginUC, authService)
	userHandler := NewUserHandler(authService)
	logoutHandler := NewLogoutHandler()

	r.Post("/api/v1/auth/register", shared.Action(regHandler.Handle))
	r.Post("/api/v1/auth/login", shared.Action(loginHandler.Handle))
	r.Post("/api/v1/auth/logout", shared.Action(logoutHandler.Handle))
	r.With(shared.AuthMiddleware(authService)).Get("/api/v1/me", shared.Action(userHandler.Handle))
}
