// @AI_GENERATED
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// LoggerMiddleware returns a Gin middleware that logs request method, path,
// status code, and duration using the provided Logger.
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"duration", duration.String(),
		)
	}
}

// @AI_GENERATED: end
