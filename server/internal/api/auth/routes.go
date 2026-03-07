package auth

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, jwtSecret string) TokenProvider {
	repo := &userRepository{db: db}
	tp := &AuthService{secret: []byte(jwtSecret)}

	// Public auth routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleRegister(w, r, repo, tp)
		}))
		r.Post("/login", shared.Action(func(w http.ResponseWriter, r *http.Request) error {
			return HandleLogin(w, r, repo, tp)
		}))
		r.Post("/logout", shared.Action(HandleLogout))
	})

	// Protected user routes
	r.With(AuthMiddleware(tp)).Route("/api/v1", func(r chi.Router) {
		r.Get("/me", shared.Action(HandleMe))
	})

	return tp
}
