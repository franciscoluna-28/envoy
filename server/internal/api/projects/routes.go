package projects

import (
	"github.com/franciscoluna/envoy/server/internal/api/auth"
	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, authService auth.TokenProvider) {
	r.With(auth.AuthMiddleware(authService)).Route("/api/v1/projects", func(r chi.Router) {
		r.Post("/", shared.Action(HandleCreateProject))
		r.Get("/", shared.Action(HandleListProjects))
		r.Get("/{id}", shared.Action(HandleGetProject))
		r.Put("/{id}", shared.Action(HandleUpdateProject))
		r.Delete("/{id}", shared.Action(HandleDeleteProject))
	})
}
