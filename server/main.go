package main

import (
	"log"
	"net/http"
	"os"

	"github.com/franciscoluna/envoy/server/internal/api/infra"
	"github.com/franciscoluna/envoy/server/internal/api/repositories"
	"github.com/franciscoluna/envoy/server/internal/api/routes"
	"github.com/franciscoluna/envoy/server/internal/api/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// TODO: Create proper env validation util
// TODO: Create proper database migration setup
// TODO: Improve app boostrap organization
// TODO: Add support for port mapping

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "./envoy.db"
	}

	db, err := sqlx.Open("sqlite", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Printf("Database created/connected at: %s", dbURL)

	// TODO: Move this to a proper sql file and add support for db migrations
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`
	_, err = db.Exec(schema)

	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	userRepo := repositories.NewUserRepository(db)

	regUC := usecases.NewRegisterUseCase(userRepo.(*repositories.UserRepository))
	loginUC := usecases.NewLoginWithEmailUseCase(userRepo.(*repositories.UserRepository))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	authService := infra.NewAuthService(jwtSecret)

	router := chi.NewRouter()

	routes.RegisterRoutes(router, authService, regUC, loginUC)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		infra.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
