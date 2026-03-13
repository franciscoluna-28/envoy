package main

import (
	"log"
	"net/http"

	_ "newserver/docs"
	"newserver/internal/auth"
	"newserver/internal/database"
	"newserver/internal/projects"
	"newserver/internal/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "modernc.org/sqlite"
)

// @title Envoy Auth API
// @version 1.0
// @description Authentication API for Envoy application
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := shared.LoadConfig()

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	v := validator.New()
	r := chi.NewRouter()
	authRepo := auth.NewRepository(db)
	jwtProvider := auth.NewJWTProvider(cfg.JWTSecret)
	authHandler := auth.NewHandler(authRepo, v, jwtProvider)
	authMiddleware := auth.AuthMiddleware(jwtProvider)

	projectRepo := projects.NewRepository(db)
	projectHandler := projects.NewHandler(projectRepo, v)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		Debug:            false,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the newborn server!"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler())

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", authHandler.Register)
				r.Post("/login", authHandler.Login)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)
				r.Get("/auth/me", authHandler.GetMe)
				r.Post("/auth/logout", authHandler.Logout)
				r.Get("/projects", projectHandler.GetAllProjects)
				r.Post("/projects", projectHandler.CreateProject)
				r.Get("/projects/{id}", projectHandler.GetProject)
				r.Put("/projects/{id}", projectHandler.UpdateProject)
				r.Delete("/projects/{id}", projectHandler.DeleteProject)
			})
		})
	})

	log.Println("Server starting on :" + cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
