package approval

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Decision represents the user's approval decision.
type Decision string

const (
	DecisionApprove Decision = "approve"
	DecisionDeny    Decision = "deny"
)

// PendingApproval holds a blocked tool call waiting for user decision.
type PendingApproval struct {
	ID        string
	SessionID string
	ToolName  string
	Args      map[string]any
	CreatedAt time.Time
	ch        chan Decision
}

// Manager manages pending tool approvals across sessions.
type Manager struct {
	mu      sync.Mutex
	pending map[string]*PendingApproval // key: approval ID
}

// NewManager creates a new approval Manager.
func NewManager() *Manager {
	return &Manager{pending: make(map[string]*PendingApproval)}
}

// RequestWithID blocks on a pre-built PendingApproval until resolved or ctx cancelled.
func (m *Manager) RequestWithID(ctx context.Context, pa *PendingApproval) (bool, string, error) {
	ch := make(chan Decision, 1)
	pa.ch = ch

	m.mu.Lock()
	m.pending[pa.ID] = pa
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.pending, pa.ID)
		m.mu.Unlock()
	}()

	select {
	case decision := <-ch:
		return decision == DecisionApprove, pa.ID, nil
	case <-ctx.Done():
		return false, pa.ID, ctx.Err()
	}
}

// Request blocks until the user approves or denies the tool call, or ctx is cancelled.
// Returns true if approved, false if denied or cancelled.
func (m *Manager) Request(ctx context.Context, sessionID, toolName string, args map[string]any) (bool, string, error) {
	id := fmt.Sprintf("appr-%d", time.Now().UnixNano())
	ch := make(chan Decision, 1)

	pa := &PendingApproval{
		ID:        id,
		SessionID: sessionID,
		ToolName:  toolName,
		Args:      args,
		CreatedAt: time.Now(),
		ch:        ch,
	}

	m.mu.Lock()
	m.pending[id] = pa
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.pending, id)
		m.mu.Unlock()
	}()

	select {
	case decision := <-ch:
		return decision == DecisionApprove, id, nil
	case <-ctx.Done():
		return false, id, ctx.Err()
	}
}

// Resolve submits a decision for the given approval ID.
// Returns an error if the approval ID is not found.
func (m *Manager) Resolve(approvalID string, decision Decision) error {
	m.mu.Lock()
	pa, ok := m.pending[approvalID]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("approval %s not found or already resolved", approvalID)
	}

	select {
	case pa.ch <- decision:
		return nil
	default:
		return fmt.Errorf("approval %s already resolved", approvalID)
	}
}

// ListPending returns all pending approvals for a session.
func (m *Manager) ListPending(sessionID string) []*PendingApproval {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []*PendingApproval
	for _, pa := range m.pending {
		if pa.SessionID == sessionID {
			result = append(result, pa)
		}
	}
	return result
}
