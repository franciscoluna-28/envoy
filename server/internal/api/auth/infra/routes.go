package infra

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/franciscoluna/envoy/server/internal/shared"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, authService auth.TokenProvider, regUC *application.RegisterWithEmailUseCase, loginUC *application.LoginWithEmailUseCase) {
	regHandler := NewRegisterHandler(regUC, authService)
	loginHandler := NewLoginHandler(loginUC, authService)

	userHandler := NewUserHandler()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", makeHandler(shared.WithBody(regHandler.Handle)))
			r.Post("/login", makeHandler(shared.WithBody(loginHandler.Handle)))
		})

		r.Group(func(r chi.Router) {
			r.Use(shared.AuthMiddleware(authService))
			r.Get("/me", makeHandler(userHandler.Me))
		})
	})
}

func makeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if appErr, ok := err.(*shared.AppError); ok {
				shared.Send(w, appErr.Status, shared.APIResponse{
					Status:  appErr.Status,
					Error:   appErr.Code,
					Message: appErr.Msg,
					Success: false,
				})
				return
			}

			shared.InternalError(w)
		}
	}
}
