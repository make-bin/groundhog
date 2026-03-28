// @AI_GENERATED
package service

import (
	"context"
	"time"
)

// AuditLog represents a single audit log entry.
type AuditLog struct {
	ID           uint
	Action       string
	PrincipalID  string
	ResourceType string
	ResourceID   string
	Details      string
	SourceIP     string
	CreatedAt    time.Time
}

// AuditFilter defines filter criteria for querying audit logs.
type AuditFilter struct {
	Action      *string
	PrincipalID *string
	Page        int
	PageSize    int
}

// AuditService defines the interface for recording and querying audit events.
type AuditService interface {
	Record(ctx context.Context, action, principalID, resourceType, resourceID, details, sourceIP string) error
	Query(ctx context.Context, filter AuditFilter) ([]*AuditLog, int, error)
}

// @AI_GENERATED: end
