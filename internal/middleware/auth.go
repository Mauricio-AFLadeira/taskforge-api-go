package middleware

import "net/http"

// RequireAuth validates JWT bearer tokens. Stub until Milestone 3 — Authentication.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
