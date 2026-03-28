package middleware

import "github.com/gin-gonic/gin"

// UserIDMiddleware extracts user identity from X-User-ID header (or query param)
// and sets it in the gin context as "user_id".
// This is a development convenience middleware; replace with JWT auth in production.
func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = c.Query("user_id")
		}
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	}
}
