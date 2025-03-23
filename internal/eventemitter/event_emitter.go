package eventemitter

import (
	"reflect"
	"sync"
)

type Handler func(any)
type Capturer func(any)

type EventEmitter struct {
	mu        sync.RWMutex
	handlers  map[string][]Handler
	capturers []Capturer
}

func New() *EventEmitter {
	return &EventEmitter{handlers: make(map[string][]Handler), capturers: make([]Capturer, 0)}
}

func (e *EventEmitter) On(event any, handler Handler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	name := reflect.TypeOf(event).String()
	e.handlers[name] = append(e.handlers[name], handler)
}

func (e *EventEmitter) AddCapturer(handler Capturer) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.capturers = append(e.capturers, handler)
}

func (e *EventEmitter) Emit(event any) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	name := reflect.TypeOf(event).String()
	for _, handler := range e.handlers[name] {
		handler(event)
	}
	for _, capturer := range e.capturers {
		capturer(event)
	}
}
