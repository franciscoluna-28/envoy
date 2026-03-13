package auth

import (
	"context"
	"net/http"
	response "newserver/internal/shared"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(tp TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Unauthorized"})
				return
			}

			claims, err := tp.ParseToken(cookie.Value)
			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid session"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims["sub"])

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
