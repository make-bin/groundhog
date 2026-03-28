package adk

import "context"

type contextKey string

const MemoryUserIDKey contextKey = "memory_user_id"

// WithMemoryUserID returns a new context with the userID stored under MemoryUserIDKey.
func WithMemoryUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, MemoryUserIDKey, userID)
}

// MemoryUserIDFromContext retrieves the userID from the context.
func MemoryUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(MemoryUserIDKey).(string)
	return userID, ok && userID != ""
}
