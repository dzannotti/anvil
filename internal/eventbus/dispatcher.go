package eventbus

import (
	"reflect"

	"github.com/zyedidia/generic/stack"
)

// EventType returns the string representation of the given event value's type
func EventType(event any) string {
	return reflect.TypeOf(event).String()
}

type GenericEventHandler func(data Event)

type Dispatcher struct {
	events           []Event
	stack            stack.Stack[Event]
	allSubscribers   []GenericEventHandler
	eventSubscribers map[string][]GenericEventHandler
}

func (s *Dispatcher) SubscribeAll(handler GenericEventHandler) {
	s.allSubscribers = append(s.allSubscribers, handler)
}

func (s *Dispatcher) Subscribe(eventType string, handler GenericEventHandler) {
	if s.eventSubscribers == nil {
		s.eventSubscribers = make(map[string][]GenericEventHandler)
	}
	s.eventSubscribers[eventType] = append(s.eventSubscribers[eventType], handler)
}

func (s *Dispatcher) Emit(event any) {
	s.Begin(event)
	s.End()
}

func (s *Dispatcher) Begin(data any) {
	eventType := reflect.TypeOf(data).String()
	event := Event{Kind: eventType, Data: data}
	event.Depth = s.stack.Size()
	s.events = append(s.events, event)
	s.stack.Push(event)
	s.emit(event)
}

func (s *Dispatcher) End() {
	if s.stack.Size() == 0 {
		return
	}
	last := s.stack.Pop()
	event := Event{Kind: last.Kind, Data: last.Data}
	event.Depth = s.stack.Size() + 1
	event.End = true
	s.events = append(s.events, event)
	s.emit(event)
}

func (s *Dispatcher) emit(event Event) {
	// Emit to all subscribers
	for _, subscriber := range s.allSubscribers {
		subscriber(event)
	}

	// Emit to event-specific subscribers
	if handlers, exists := s.eventSubscribers[event.Kind]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
