package auth

import (
	"context"
	"net/http"

	"github.com/franciscoluna/envoy/server/internal/shared"
)

func AuthMiddleware(tp TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(AuthCookieName)
			if err != nil {
				shared.BadRequest(w, shared.ErrUnauthorized, "Missing or invalid session cookie", nil)
				return
			}
			claims, err := tp.ParseToken(cookie.Value)
			if err != nil {
				shared.BadRequest(w, shared.ErrUnauthorized, "Session expired or invalid", nil)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", claims["sub"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
