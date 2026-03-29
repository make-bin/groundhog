package router

import (
	"github.com/gin-gonic/gin"
	"github.com/make-bin/groundhog/pkg/interface/http/handler"
)

// RegisterSessionRoutes registers session-related routes on the given router group.
func RegisterSessionRoutes(rg *gin.RouterGroup, sessionHandler handler.SessionHandler) {
	sessions := rg.Group("/sessions")
	{
		sessions.GET("", sessionHandler.List)
		sessions.POST("", sessionHandler.Create)
		sessions.GET("/:id", sessionHandler.Get)
		sessions.DELETE("/:id", sessionHandler.Delete)
		sessions.POST("/:id/messages", sessionHandler.SendMessage)
		sessions.POST("/:id/messages/stream", sessionHandler.StreamMessage)
		sessions.GET("/:id/approvals", sessionHandler.ListApprovals)
		sessions.POST("/:id/approvals/:approval_id", sessionHandler.ResolveApproval)
	}
}

// RegisterChannelRoutes registers channel-related routes on the given router group.
func RegisterChannelRoutes(rg *gin.RouterGroup, channelHandler handler.ChannelHandler) {
	channels := rg.Group("/channels")
	{
		channels.GET("", channelHandler.List)
		channels.POST("", channelHandler.Create)
		channels.DELETE("/:id", channelHandler.Delete)
		channels.GET("/:id/status", channelHandler.Status)
	}
}

// RegisterSecurityRoutes registers security-related routes on the given router group.
func RegisterSecurityRoutes(rg *gin.RouterGroup, securityHandler handler.SecurityHandler) {
	security := rg.Group("/security")
	{
		security.GET("/audit", securityHandler.AuditLogs)
	}
}

// RegisterMemoryRoutes registers memory-related routes on the given router group.
func RegisterMemoryRoutes(rg *gin.RouterGroup, memoryHandler handler.MemoryHandler) {
	memories := rg.Group("/memories")
	{
		memories.GET("", memoryHandler.List)
		memories.POST("", memoryHandler.Create)
		memories.GET("/:id", memoryHandler.Get)
		memories.PUT("/:id", memoryHandler.Update)
		memories.DELETE("/:id", memoryHandler.Delete)
		memories.POST("/search", memoryHandler.Search)
	}
}

// RegisterAgentRoutes registers agent-related routes on the given router group.
func RegisterAgentRoutes(rg *gin.RouterGroup, agentHandler handler.AgentHandler) {
	agents := rg.Group("/agents")
	{
		agents.GET("", agentHandler.List)
	}
}
