package events

import (
	"reflect"
	"sync"
)

var registry = &BusRegistry{
	buses: make(map[reflect.Type]any),
}

type BusRegistry struct {
	mu    sync.RWMutex
	buses map[reflect.Type]any
}

// Generic function, not a method
func GetBus[T any]() *EventBus[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	registry.mu.RLock()
	bus, ok := registry.buses[typ]
	registry.mu.RUnlock()

	if ok {
		return bus.(*EventBus[T])
	}

	// Double-checked locking
	registry.mu.Lock()
	defer registry.mu.Unlock()

	// Someone else might have created it while we waited
	if bus, ok := registry.buses[typ]; ok {
		return bus.(*EventBus[T])
	}

	newBus := NewEventBus[T]()
	registry.buses[typ] = newBus
	return newBus
}

// EventBus[T] handles subscriptions and publishing for a specific event type T
type EventBus[T any] struct {
	subscribers []func(T)
	mu          sync.RWMutex
}

func NewEventBus[T any]() *EventBus[T] {
	return &EventBus[T]{}
}

func (b *EventBus[T]) Subscribe(fn func(T)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, fn)
}

func (b *EventBus[T]) Publish(event T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.subscribers {
		go sub(event) // run handlers in goroutines
	}
}
