// @AI_GENERATED
package entity

import "github.com/make-bin/groundhog/pkg/domain/messaging/vo"

// Route represents a routing rule that maps an account binding to a session.
type Route struct {
	id             string
	accountBinding vo.AccountBinding
	priority       int
}

// NewRoute creates a new Route.
func NewRoute(id string, binding vo.AccountBinding, priority int) *Route {
	return &Route{
		id:             id,
		accountBinding: binding,
		priority:       priority,
	}
}

// ID returns the route identifier.
func (r *Route) ID() string { return r.id }

// AccountBinding returns the account binding for this route.
func (r *Route) AccountBinding() vo.AccountBinding { return r.accountBinding }

// Priority returns the routing priority (higher value = higher priority).
func (r *Route) Priority() int { return r.priority }

// @AI_GENERATED: end
