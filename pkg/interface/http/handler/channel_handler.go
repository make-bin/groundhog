// @AI_GENERATED
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/application/dto"
	appservice "github.com/make-bin/groundhog/pkg/application/service"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// ChannelHandler defines the HTTP handler interface for channel management.
type ChannelHandler interface {
	List(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Status(c *gin.Context)
}

type channelHandler struct {
	ChannelAppService appservice.ChannelAppService `inject:""`
}

// NewChannelHandler creates a new ChannelHandler. Dependencies are injected via struct tags.
func NewChannelHandler() ChannelHandler {
	return &channelHandler{}
}

// List handles GET /api/v1/channels
func (h *channelHandler) List(c *gin.Context) {
	channels, err := h.ChannelAppService.ListChannels(c.Request.Context())
	if err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	response.Success(c, channels)
}

// Create handles POST /api/v1/channels
func (h *channelHandler) Create(c *gin.Context) {
	var req dto.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	ch, err := h.ChannelAppService.CreateChannel(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 201, "data": ch})
}

// Delete handles DELETE /api/v1/channels/:id
func (h *channelHandler) Delete(c *gin.Context) {
	id, err := vo.NewChannelID(c.Param("id"))
	if err != nil {
		response.Error(c, bcode.ErrInvalidRequest, err)
		return
	}
	if err := h.ChannelAppService.DeleteChannel(c.Request.Context(), id); err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Status handles GET /api/v1/channels/:id/status
func (h *channelHandler) Status(c *gin.Context) {
	id, err := vo.NewChannelID(c.Param("id"))
	if err != nil {
		response.Error(c, bcode.ErrInvalidRequest, err)
		return
	}
	ch, err := h.ChannelAppService.GetChannelStatus(c.Request.Context(), id)
	if err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	response.Success(c, ch)
}

// @AI_GENERATED: end
