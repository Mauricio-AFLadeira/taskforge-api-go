package middleware

import (
	"net/http"
	"strings"

	"github.com/mauricio-reportei/taskforge-api-go/internal/auth"
	"github.com/mauricio-reportei/taskforge-api-go/internal/shared"
)

const bearerPrefix = "Bearer "

// RequireAuth validates JWT bearer tokens and attaches user id and email to the request context.
func RequireAuth(jwtSecret string, next http.Handler) http.Handler {
	secret := []byte(jwtSecret)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("Authorization")
		if !strings.HasPrefix(raw, bearerPrefix) {
			_ = shared.WriteError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(raw, bearerPrefix))
		if tokenStr == "" {
			_ = shared.WriteError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		userID, email, err := auth.ParseAccess(secret, tokenStr)
		if err != nil {
			_ = shared.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		ctx := shared.WithEmail(shared.WithUserID(r.Context(), userID), email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
