// @AI_GENERATED
package eventbus

import (
	"reflect"
	"sync"
	"testing"
)

type testEvent struct {
	Name string
}

func TestNewEventBus(t *testing.T) {
	bus := NewEventBus()
	if bus == nil {
		t.Fatal("NewEventBus() returned nil")
	}
}

func TestSubscribeAndPublish(t *testing.T) {
	bus := NewEventBus()
	var received interface{}

	eventType := reflect.TypeOf(testEvent{}).String()
	bus.Subscribe(eventType, func(event interface{}) {
		received = event
	})

	bus.Publish(testEvent{Name: "hello"})

	if received == nil {
		t.Fatal("handler was not called")
	}
	e, ok := received.(testEvent)
	if !ok {
		t.Fatal("received event has wrong type")
	}
	if e.Name != "hello" {
		t.Fatalf("expected Name=hello, got %s", e.Name)
	}
}

func TestMultipleSubscribers(t *testing.T) {
	bus := NewEventBus()
	var count int

	eventType := reflect.TypeOf(testEvent{}).String()
	bus.Subscribe(eventType, func(_ interface{}) { count++ })
	bus.Subscribe(eventType, func(_ interface{}) { count++ })
	bus.Subscribe(eventType, func(_ interface{}) { count++ })

	bus.Publish(testEvent{Name: "multi"})

	if count != 3 {
		t.Fatalf("expected 3 handler calls, got %d", count)
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	bus := NewEventBus()
	// Should not panic when publishing with no subscribers.
	bus.Publish(testEvent{Name: "nobody"})
}

func TestConcurrentSubscribeAndPublish(t *testing.T) {
	bus := NewEventBus()
	eventType := reflect.TypeOf(testEvent{}).String()

	var wg sync.WaitGroup
	var mu sync.Mutex
	callCount := 0

	// Subscribe concurrently.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.Subscribe(eventType, func(_ interface{}) {
				mu.Lock()
				callCount++
				mu.Unlock()
			})
		}()
	}
	wg.Wait()

	// Publish and verify all 10 handlers are called.
	bus.Publish(testEvent{Name: "concurrent"})

	if callCount != 10 {
		t.Fatalf("expected 10 handler calls, got %d", callCount)
	}
}

// @AI_GENERATED: end
