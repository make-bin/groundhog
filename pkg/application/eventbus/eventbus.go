// @AI_GENERATED
package eventbus

import (
	"reflect"
	"sync"
)

// EventBus defines the interface for an in-process domain event bus.
type EventBus interface {
	Publish(event interface{})
	Subscribe(eventType string, handler func(event interface{}))
}

// eventBus is a synchronous in-process implementation of EventBus.
type eventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]func(event interface{})
}

// NewEventBus creates a new EventBus instance.
func NewEventBus() EventBus {
	return &eventBus{
		subscribers: make(map[string][]func(event interface{})),
	}
}

// Subscribe registers a handler for the given event type.
// The same event type supports multiple subscribers.
func (b *eventBus) Subscribe(eventType string, handler func(event interface{})) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

// Publish synchronously invokes all registered handlers whose event type
// matches the runtime type name of the given event.
func (b *eventBus) Publish(event interface{}) {
	eventType := reflect.TypeOf(event).String()

	b.mu.RLock()
	handlers := make([]func(event interface{}), len(b.subscribers[eventType]))
	copy(handlers, b.subscribers[eventType])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}

// @AI_GENERATED: end
