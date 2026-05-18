package shared

import "context"

type ctxKey int

const (
	keyUserID ctxKey = iota + 1
	keyEmail
)

// WithUserID returns a child context carrying the authenticated user id.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

// WithEmail returns a child context carrying the authenticated user's email.
func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, keyEmail, email)
}

// UserIDFromContext reads the user id set by auth middleware.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyUserID).(string)
	return v, ok && v != ""
}

// EmailFromContext reads the email set by auth middleware.
func EmailFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyEmail).(string)
	return v, ok && v != ""
}
