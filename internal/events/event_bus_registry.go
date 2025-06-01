package events

import (
	"reflect"
	"sync"
)

var registry = &busRegistry{
	buses: make(map[reflect.Type]any),
}

type busRegistry struct {
	mu    sync.RWMutex
	buses map[reflect.Type]any
}

func getBus[T any]() *eventBus[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	registry.mu.RLock()
	bus, ok := registry.buses[typ]
	registry.mu.RUnlock()

	if ok {
		return bus.(*eventBus[T])
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if bus, ok := registry.buses[typ]; ok {
		return bus.(*eventBus[T])
	}

	newBus := newEventBus[T]()
	registry.buses[typ] = newBus
	return newBus
}

type eventBus[T any] struct {
	subscribers []func(T)
	mu          sync.RWMutex
	last        T
}

func newEventBus[T any]() *eventBus[T] {
	return &eventBus[T]{}
}

func (b *eventBus[T]) subscribe(fn func(T)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, fn)
}

func (b *eventBus[T]) publish(event T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.last = event
	for _, sub := range b.subscribers {
		go sub(event)
	}
}

func (b *eventBus[T]) current() T {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.last
}

func Subscribe[T any](fn func(T)) {
	getBus[T]().subscribe(fn)
}

func Publish[T any](event T) {
	getBus[T]().publish(event)
}

func Current[T any]() T {
	return getBus[T]().current()
}
