package eventbus

import "github.com/zyedidia/generic/stack"

type EventHandler func(data Event)

type Dispatcher struct {
	events      []Event
	stack       stack.Stack[Event]
	subscribers []EventHandler
}

func (s *Dispatcher) Subscribe(handler EventHandler) {
	s.subscribers = append(s.subscribers, handler)
}

func (s *Dispatcher) Add(kind string, data any) {
	s.Start(kind, data)
	s.End()
}

func (s *Dispatcher) Start(kind string, data any) {
	event := Event{Kind: kind, Data: data}
	event.Depth = s.stack.Size()
	s.events = append(s.events, event)
	s.stack.Push(event)
	s.Emit(event)
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
	s.Emit(event)
}

func (s *Dispatcher) Emit(event Event) {
	for _, subscriber := range s.subscribers {
		subscriber(event)
	}
}
