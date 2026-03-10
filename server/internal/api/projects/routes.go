package projects

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/api/auth"
	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, authService auth.TokenProvider) {
	repo := &ProjectRepo{db: db}

	r.With(auth.AuthMiddleware(authService)).Route("/api/v1/projects", func(r chi.Router) {
		r.Post("/", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleCreateProject(w, r, repo)
		}))
		r.Get("/", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleListProjects(w, r, repo)
		}))
		r.Get("/{id}", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleGetProject(w, r, repo)
		}))
		r.Put("/{id}", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleUpdateProject(w, r, repo)
		}))
		r.Delete("/{id}", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleDeleteProject(w, r, repo)
		}))
	})
}
