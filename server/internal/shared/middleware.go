package shared

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/cors"
)

func WithBody[T any](h func(w http.ResponseWriter, r *http.Request, body T) error) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var body T
		if err := Decode(r, &body); err != nil {
			return NewAppError(http.StatusBadRequest, ErrInvalidInput, fmt.Sprintf("Invalid request body: %v", err))
		}
		if err := Validate(body); err != nil {
			return NewAppError(http.StatusBadRequest, ErrInvalidInput, fmt.Sprintf("Validation failed: %v", err))
		}
		return h(w, r, body)
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("[%s] %s %s | %v\n",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("PANIC RECOVERED: %v\n", err)

				InternalError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CORS() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		Debug:            false,
	})

	return c.Handler
}
