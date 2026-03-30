package auth

import (
	"context"
	"net/http"
	response "newserver/internal/shared"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// GetUserIDFromContext safely extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func AuthMiddleware(tp TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("auth_token")
			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Unauthorized"})
				return
			}

			claims, err := tp.ParseToken(cookie.Value)
			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid session"})
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid user ID in token"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
