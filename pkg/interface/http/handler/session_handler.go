package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/make-bin/groundhog/pkg/application/dto"
	appservice "github.com/make-bin/groundhog/pkg/application/service"
	conversation "github.com/make-bin/groundhog/pkg/domain/conversation"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// SessionHandler defines the HTTP handler interface for session management.
type SessionHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	SendMessage(c *gin.Context)
	StreamMessage(c *gin.Context)
	ResolveApproval(c *gin.Context)
	ListApprovals(c *gin.Context)
}

type sessionHandler struct {
	AgentAppService appservice.AgentAppService `inject:""`
}

// NewSessionHandler creates a new SessionHandler. Dependencies are injected via struct tags.
func NewSessionHandler() SessionHandler {
	return &sessionHandler{}
}

// Create handles POST /api/v1/sessions
func (h *sessionHandler) Create(c *gin.Context) {
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	sess, err := h.AgentAppService.CreateSession(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	response.Success(c, sess)
}

// List handles GET /api/v1/sessions
func (h *sessionHandler) List(c *gin.Context) {
	var req dto.SessionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	result, err := h.AgentAppService.ListSessions(c.Request.Context(), req)
	if err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	response.Success(c, result)
}

// Get handles GET /api/v1/sessions/:id
func (h *sessionHandler) Get(c *gin.Context) {
	id, err := vo.NewSessionID(c.Param("id"))
	if err != nil {
		response.Error(c, bcode.ErrInvalidRequest, err)
		return
	}
	sess, err := h.AgentAppService.GetSession(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, conversation.ErrSessionNotFound) {
			response.Error(c, bcode.ErrNotFound, err)
			return
		}
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	response.Success(c, sess)
}

// Delete handles DELETE /api/v1/sessions/:id
func (h *sessionHandler) Delete(c *gin.Context) {
	id, err := vo.NewSessionID(c.Param("id"))
	if err != nil {
		response.Error(c, bcode.ErrInvalidRequest, err)
		return
	}
	if err := h.AgentAppService.DeleteSession(c.Request.Context(), id); err != nil {
		if errors.Is(err, conversation.ErrSessionNotFound) {
			response.Error(c, bcode.ErrNotFound, err)
			return
		}
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// SendMessage handles POST /api/v1/sessions/:id/messages
func (h *sessionHandler) SendMessage(c *gin.Context) {
	id, err := vo.NewSessionID(c.Param("id"))
	if err != nil {
		response.Error(c, bcode.ErrInvalidRequest, err)
		return
	}
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	turn, err := h.AgentAppService.ExecuteTurn(c.Request.Context(), id, req.UserInput)
	if err != nil {
		if errors.Is(err, conversation.ErrSessionNotFound) {
			response.Error(c, bcode.ErrNotFound, err)
			return
		}
		response.Error(c, bcode.ErrInternal, err)
		return
	}
	response.Success(c, turn)
}

// StreamMessage handles POST /api/v1/sessions/:id/messages/stream
// Streams the LLM response as Server-Sent Events (SSE).
// Event format: data: <json StreamEvent>\n\n
func (h *sessionHandler) StreamMessage(c *gin.Context) {
	id, err := vo.NewSessionID(c.Param("id"))
	if err != nil {
		response.Error(c, bcode.ErrInvalidRequest, err)
		return
	}
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}

	eventCh := h.AgentAppService.StreamTurn(c.Request.Context(), id, req.UserInput)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, canFlush := c.Writer.(http.Flusher)

	for event := range eventCh {
		data, _ := json.Marshal(event)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		if canFlush {
			flusher.Flush()
		}
		if event.Type == "done" || event.Type == "error" {
			return
		}
	}
}

// ResolveApproval handles POST /api/v1/sessions/:id/approvals/:approval_id
func (h *sessionHandler) ResolveApproval(c *gin.Context) {
	approvalID := c.Param("approval_id")
	var req dto.ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, bcode.ErrValidationFailed, err)
		return
	}
	if err := h.AgentAppService.ResolveApproval(approvalID, req.Decision); err != nil {
		response.Error(c, bcode.ErrNotFound, err)
		return
	}
	response.Success(c, nil)
}

// ListApprovals handles GET /api/v1/sessions/:id/approvals
// Returns all pending tool approvals for the session.
func (h *sessionHandler) ListApprovals(c *gin.Context) {
	sessionID := c.Param("id")
	approvals := h.AgentAppService.ListApprovals(sessionID)
	response.Success(c, approvals)
}
