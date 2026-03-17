package auth

import (
	"context"
	"fmt"
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
			fmt.Printf("AuthMiddleware: Processing request to %s\n", r.URL.Path)

			cookie, err := r.Cookie("auth_token")
			if err != nil {
				fmt.Printf("AuthMiddleware: Cookie error - %v\n", err)
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Unauthorized"})
				return
			}

			fmt.Printf("AuthMiddleware: Found cookie, value length: %d\n", len(cookie.Value))

			claims, err := tp.ParseToken(cookie.Value)
			if err != nil {
				fmt.Printf("AuthMiddleware: Token parse error - %v\n", err)
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid session"})
				return
			}

			fmt.Printf("AuthMiddleware: Token parsed successfully, claims: %+v\n", claims)

			userID, ok := claims["sub"].(string)
			if !ok {
				fmt.Printf("AuthMiddleware: Invalid subject claim\n")
				response.WriteJSON(w, http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid user ID in token"})
				return
			}

			fmt.Printf("AuthMiddleware: User ID extracted: %s\n", userID)

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			fmt.Printf("AuthMiddleware: Context updated, proceeding to next handler\n")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
