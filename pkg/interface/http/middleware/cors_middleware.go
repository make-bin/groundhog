// @AI_GENERATED
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a Gin middleware that handles Cross-Origin Resource Sharing.
// Currently allows all origins for development; restrict in production as needed.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// @AI_GENERATED: end
