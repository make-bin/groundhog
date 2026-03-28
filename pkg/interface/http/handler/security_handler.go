// @AI_GENERATED
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	identitysvc "github.com/make-bin/groundhog/pkg/domain/identity/service"
	"github.com/make-bin/groundhog/pkg/interface/http/response"
	"github.com/make-bin/groundhog/pkg/utils/bcode"
)

// SecurityHandler defines the HTTP handler interface for security operations.
type SecurityHandler interface {
	AuditLogs(c *gin.Context)
}

type securityHandler struct {
	AuditService identitysvc.AuditService `inject:""`
}

// NewSecurityHandler creates a new SecurityHandler.
func NewSecurityHandler() SecurityHandler {
	return &securityHandler{}
}

// AuditLogs handles GET /api/v1/security/audit
func (h *securityHandler) AuditLogs(c *gin.Context) {
	filter := identitysvc.AuditFilter{}

	if action := c.Query("action"); action != "" {
		filter.Action = &action
	}
	if principalID := c.Query("principal_id"); principalID != "" {
		filter.PrincipalID = &principalID
	}
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = p
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil {
			filter.PageSize = ps
		}
	}

	logs, total, err := h.AuditService.Query(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, bcode.ErrInternal, err)
		return
	}

	response.Success(c, gin.H{
		"items": logs,
		"total": total,
		"page":  filter.Page,
	})
}

// @AI_GENERATED: end
