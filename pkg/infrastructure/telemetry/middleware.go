// @AI_GENERATED
package telemetry

import (
	"time"

	"github.com/gin-gonic/gin"
)

// OtelMiddleware returns a Gin middleware that records HTTP request metrics.
// It increments the request counter and accumulates request duration.
func OtelMiddleware(m *Metrics) gin.HandlerFunc {
	if m == nil {
		// No-op if telemetry is not initialised.
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		m.RecordHTTPRequest(time.Since(start))
	}
}

// @AI_GENERATED: end
