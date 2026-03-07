package environments

import (
	"github.com/franciscoluna/envoy/server/internal/api/auth"
	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, authService auth.TokenProvider) {
	r.With(auth.AuthMiddleware(authService)).Route("/api/v1/environments", func(r chi.Router) {
		r.Post("/", shared.Action(HandleCreateEnvironment))
	})

	r.Post("/api/v1/database/test-connection", shared.Action(HandleTestConnection))
}
