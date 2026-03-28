// @AI_GENERATED
package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/infrastructure/service"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// from the Authorization header and sets the principal claims in context.
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, bcode.ErrUnauthorized, errors.New("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, bcode.ErrUnauthorized, errors.New("invalid authorization format"))
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			response.Error(c, bcode.ErrUnauthorized, errors.New("empty token"))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				response.Error(c, bcode.ErrTokenExpired, err)
				c.Abort()
				return
			}
			if errors.Is(err, service.ErrTokenInvalid) {
				response.Error(c, bcode.ErrTokenInvalid, err)
				c.Abort()
				return
			}
			response.Error(c, bcode.ErrUnauthorized, err)
			c.Abort()
			return
		}

		c.Set("principal", claims)
		c.Next()
	}
}

// @AI_GENERATED: end
