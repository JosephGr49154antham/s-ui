package middleware

import (
	"net/http"
	"strings"

	"github.com/s-ui/auth"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix       = "Bearer "
)

// Authenticate is an HTTP middleware that validates JWT tokens.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(AuthorizationHeader)
		if header == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(header, BearerPrefix) {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, BearerPrefix)
		claims, err := auth.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		_ = claims
		next.ServeHTTP(w, r)
	})
}
