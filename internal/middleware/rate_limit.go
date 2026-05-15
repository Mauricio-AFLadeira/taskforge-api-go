package middleware

import "net/http"

// RateLimiter applies Redis-backed rate limiting. Stub until Milestone 8 — Redis Features.
func RateLimiter(next http.Handler) http.Handler {
	return next
}
