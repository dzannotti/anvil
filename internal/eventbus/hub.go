package eventbus

import "github.com/zyedidia/generic/stack"

type EventHandler func(data Event)

type Dispatcher struct {
	events           []Event
	stack            stack.Stack[Event]
	allSubscribers   []EventHandler
	eventSubscribers map[string][]EventHandler
}

func (s *Dispatcher) SubscribeAll(handler EventHandler) {
	s.allSubscribers = append(s.allSubscribers, handler)
}

func (s *Dispatcher) Subscribe(eventType string, handler EventHandler) {
	if s.eventSubscribers == nil {
		s.eventSubscribers = make(map[string][]EventHandler)
	}
	s.eventSubscribers[eventType] = append(s.eventSubscribers[eventType], handler)
}

func (s *Dispatcher) Emit(kind string, data any) {
	s.Begin(kind, data)
	s.End()
}

func (s *Dispatcher) Begin(kind string, data any) {
	event := Event{Kind: kind, Data: data}
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
