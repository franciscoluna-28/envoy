package main

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

	"github.com/franciscoluna/envoy/server/internal/api/auth/application"
	auth "github.com/franciscoluna/envoy/server/internal/api/auth/domain"
	"github.com/franciscoluna/envoy/server/internal/api/auth/infra"
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

type App struct {
	config      *Config
	db          *sqlx.DB
	authService auth.TokenProvider
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

func NewApp(config *Config) (*App, error) {
	db, err := sqlx.Open("sqlite", config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	authService := infra.NewAuthService(config.JWTSecret)

	return &App{
		config:      config,
		db:          db,
		authService: authService,
	}, nil
}

func runMigrations(db *sqlx.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

func (a *App) setupDependencies() (*application.RegisterWithEmailUseCase, *application.LoginWithEmailUseCase) {
	userRepo := infra.NewUserRepository(a.db)
	regUC := application.NewRegisterUseCase(userRepo)
	loginUC := application.NewLoginWithEmailUseCase(userRepo)
	return regUC, loginUC
}

func (a *App) setupRouter() http.Handler {
	regUC, loginUC := a.setupDependencies()
	router := chi.NewRouter()

	router.Use(shared.Recoverer)
	router.Use(shared.Logger)
	router.Use(shared.CORS())

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		shared.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	infra.RegisterRoutes(router, a.authService, regUC, loginUC)

	return router
}

func (a *App) run() error {
	router := a.setupRouter()

	srv := &http.Server{
		Addr:         ":" + a.config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", a.config.Port)
	log.Printf("Database connected at: %s", a.config.DatabaseURL)

	return srv.ListenAndServe()
}

func (a *App) shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	done := make(chan error, 1)
	go func() {
		done <- a.db.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Error closing database: %v", err)
			return err
		}
		log.Println("Database closed successfully")
	case <-ctx.Done():
		log.Println("Shutdown timeout reached")
		return ctx.Err()
	}

	return nil
}

// @title           Envoy API
// @version         1.0
// @description     Envoy API.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	config := loadConfig()
	app, err := NewApp(config)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := app.run(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	<-ctx.Done()
	if err := app.shutdown(context.Background()); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Server stopped")
}
