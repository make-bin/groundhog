package service

import (
	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
)

// RoutingService resolves the target session for an inbound message.
type RoutingService interface {
	Resolve(msg *inbound_message.InboundMessage) (string, error)
}
