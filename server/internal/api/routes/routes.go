package routes

import (
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/api/handlers"
	"github.com/franciscoluna/envoy/server/internal/api/infra"
	"github.com/franciscoluna/envoy/server/internal/api/usecases"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, authService *infra.AuthService, regUC *usecases.RegisterWithEmailUseCase, loginUC *usecases.LoginWithEmailUseCase) {
	regHandler := handlers.NewRegisterHandler(regUC, authService)
	loginHandler := handlers.NewLoginHandler(loginUC, authService)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", makeHandler(regHandler.Handle))
			r.Post("/login", makeHandler(loginHandler.Handle))
		})
	})
}

func makeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			infra.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}
