package main

import (
	"log"
	"net/http"

	_ "newserver/docs"
	"newserver/internal/auth"
	"newserver/internal/database"

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
	db, err := database.New("envoy.db")
	if err != nil {
		panic(err)
	}

	v := validator.New()
	r := chi.NewRouter()
	authMiddleware := auth.AuthMiddleware(auth.NewJWTProvider("secret"))

	authRepo := auth.NewRepository(db)
	authHandler := auth.NewHandler(authRepo, v)

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
			})
		})
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
