package main

// @title          Envoy API
// @description    Envoy API for managing database environments and projects
// @version        1.0.0
// @contact.name   API Support
// @host           localhost:8080
// @BasePath       /
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name session_token
// @description Authentication cookie containing JWT token

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/franciscoluna/envoy/server/docs"
	"github.com/franciscoluna/envoy/server/internal/api/auth"
	"github.com/franciscoluna/envoy/server/internal/api/environments"
	"github.com/franciscoluna/envoy/server/internal/shared"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func loadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}
	return &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "./envoy.db"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "default-secret-change-in-production"),
		Port:        getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func openDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,	
			created_by TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS environments (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			project_id TEXT NOT NULL,
			connection_string_encrypted BLOB NOT NULL,
			ssl_mode TEXT NOT NULL,
			certificates_json TEXT, 
			connection_status TEXT NOT NULL DEFAULT 'pending',
			connection_error TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,	
    		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE			
		);
	`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create core tables: %w", err)
	}
	return nil
}

func setupRouter(db *sqlx.DB, config *Config) http.Handler {
	r := chi.NewRouter()

	r.Use(shared.Recoverer)
	r.Use(shared.Logger)
	r.Use(shared.CORS())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		shared.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	authService := auth.RegisterRoutes(r, db, config.JWTSecret)
	environments.RegisterRoutes(r, authService)

	return r
}

func main() {
	config := loadConfig()

	db, err := openDB(config.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	router := setupRouter(db, config)
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", config.Port)
	log.Printf("Database connected at: %s", config.DatabaseURL)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("Server stopped")
}
