// @AI_GENERATED
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// Response represents the unified API response format.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a 200 OK response with the standard success format.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    data,
	})
}

// Error sends an error response using the given BCode and error.
func Error(c *gin.Context, bc bcode.BCode, err error) {
	c.JSON(bc.HTTPStatus, Response{
		Code:    bc.Code,
		Message: bc.Message,
		Error:   err.Error(),
	})
}

// @AI_GENERATED: end
