package handler

import (
	"github.com/gin-gonic/gin"
	appservice "github.com/make-bin/groundhog/pkg/application/service"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
)

// AgentHandler defines the HTTP handler interface for agent management.
type AgentHandler interface {
	List(c *gin.Context)
}

type agentHandler struct {
	AgentAppService appservice.AgentAppService `inject:""`
}

// NewAgentHandler creates a new AgentHandler. Dependencies are injected via struct tags.
func NewAgentHandler() AgentHandler {
	return &agentHandler{}
}

// List handles GET /api/v1/agents
func (h *agentHandler) List(c *gin.Context) {
	agents := h.AgentAppService.ListAgents()
	response.Success(c, agents)
}
