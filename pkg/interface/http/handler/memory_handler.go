// @AI_GENERATED
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/application/dto"
	appservice "github.com/make-bin/groundhog/pkg/application/service"
	memorydomain "github.com/make-bin/groundhog/pkg/domain/memory"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// MemoryHandler defines the HTTP handler interface for memory management.
type MemoryHandler interface {
	List(c *gin.Context)
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Search(c *gin.Context)
}

type memoryHandler struct {
	MemorySvc appservice.MemoryAppService `inject:""`
}

// NewMemoryHandler creates a new MemoryHandler. Dependencies are injected via struct tags.
func NewMemoryHandler() MemoryHandler {
	return &memoryHandler{}
}

func (h *memoryHandler) getUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	return uid
}

func (h *memoryHandler) handleMemoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, memorydomain.ErrMemoryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
	case errors.Is(err, memorydomain.ErrMemoryAccessDenied):
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
	case errors.Is(err, memorydomain.ErrEmptyContent), errors.Is(err, memorydomain.ErrMissingUserID):
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
	default:
		response.Error(c, bcode.ErrInternal, err)
	}
}

// List handles GET /api/v1/memories
func (h *memoryHandler) List(c *gin.Context) {
	userID := h.getUserID(c)
	offset := 0
	limit := 20
	memories, err := h.MemorySvc.ListMemories(c.Request.Context(), userID, offset, limit)
	if err != nil {
		h.handleMemoryError(c, err)
		return
	}
	response.Success(c, memories)
}

// Create handles POST /api/v1/memories
func (h *memoryHandler) Create(c *gin.Context) {
	var req dto.CreateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	userID := h.getUserID(c)
	mem, err := h.MemorySvc.SaveMemory(c.Request.Context(), userID, req.Content)
	if err != nil {
		h.handleMemoryError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 201, "data": mem})
}

// Get handles GET /api/v1/memories/:id
func (h *memoryHandler) Get(c *gin.Context) {
	id := c.Param("id")
	userID := h.getUserID(c)
	mem, err := h.MemorySvc.GetMemory(c.Request.Context(), id, userID)
	if err != nil {
		h.handleMemoryError(c, err)
		return
	}
	response.Success(c, mem)
}

// Update handles PUT /api/v1/memories/:id
func (h *memoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	userID := h.getUserID(c)
	mem, err := h.MemorySvc.UpdateMemory(c.Request.Context(), id, userID, req.Content)
	if err != nil {
		h.handleMemoryError(c, err)
		return
	}
	response.Success(c, mem)
}

// Delete handles DELETE /api/v1/memories/:id
func (h *memoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := h.getUserID(c)
	if err := h.MemorySvc.DeleteMemory(c.Request.Context(), id, userID); err != nil {
		h.handleMemoryError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Search handles POST /api/v1/memories/search
func (h *memoryHandler) Search(c *gin.Context) {
	var req dto.MemorySearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	userID := h.getUserID(c)
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	results, err := h.MemorySvc.SearchMemory(c.Request.Context(), userID, req.Query, limit)
	if err != nil {
		h.handleMemoryError(c, err)
		return
	}
	response.Success(c, results)
}

// @AI_GENERATED: end
