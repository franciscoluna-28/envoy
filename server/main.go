package main

import (
	"log"
	"net/http"
	"strings"

	_ "newserver/docs"
	"newserver/internal/auth"
	"newserver/internal/database"
	"newserver/internal/environments"
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
		log.Fatalf("Failed to connect to database: %v", err)
	}

	v := validator.New()
	r := chi.NewRouter()
	authRepo := auth.NewRepository(db)
	jwtProvider := auth.NewJWTProvider(cfg.JWTSecret)
	authHandler := auth.NewHandler(authRepo, v, jwtProvider)
	authMiddleware := auth.AuthMiddleware(jwtProvider)
	postgresValidator := environments.NewPostgresValidator()

	projectRepo := projects.NewRepository(db)
	projectHandler := projects.NewHandler(projectRepo, v)

	environmentRepo := environments.NewRepository(db)
	masterKey := []byte(cfg.EncryptionKey)
	checksumKey := []byte(cfg.ChecksumKey)

	environmentHandler := environments.NewHandler(environmentRepo, v, masterKey, checksumKey, postgresValidator)

	// Validate CORS origins for security
	origins := strings.Split(cfg.AllowedOrigins, ",")
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "*" {
			log.Fatal("CORS wildcard '*' is not allowed when AllowCredentials is true")
		}
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		Debug:            cfg.Env == "development",
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		shared.WriteJSON(w, http.StatusNotFound, shared.ErrorResponse{Message: "Route not found"})
	})

	// API documentation
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
				r.Post("/projects/{id}/environments", environmentHandler.CreateEnvironment)
				r.Get("/projects/{id}/environments", environmentHandler.GetAllEnvironmentsByProjectID)
				r.Get("/environments/{id}", environmentHandler.GetEnvironmentByID)
				r.Put("/environments/{id}", environmentHandler.UpdateEnvironment)
				r.Delete("/environments/{id}", environmentHandler.DeleteEnvironment)
				r.Get("/environments/{id}/schema", environmentHandler.GetEnvironmentSchema)
				r.Post("/environments/{id}/migrations/preview", environmentHandler.PreviewEnvironmentSchemaChanges)
				r.Post("/environments/{id}/migrations", environmentHandler.RunDatabaseMigration)
				r.Get("/environments/{id}/migrations", environmentHandler.GetEnvironmentMigrations)
				r.Get("/migrations/{id}", environmentHandler.GetEnvironmentMigrationByID)
				r.Post("/environments/{id}/validate", environmentHandler.ValidateEnvironmentConnection)
				r.Post("/environments/{id}/verify-permissions", environmentHandler.VerifyDatabasePermissions)
				r.Post("/environments/{id}/test-permissions-preview", environmentHandler.TestPermissionsWithPreview)
				r.Post("/environments/{id}/test-permissions-current", environmentHandler.TestPermissionsWithCurrentSchema)
			})
		})
	})

	log.Println("Server starting on :" + cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
