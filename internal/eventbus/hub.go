package eventbus

import "github.com/zyedidia/generic/stack"

type MessageHandler func(data Message)

type Hub struct {
	events      []Message
	stack       stack.Stack[Message]
	subscribers []MessageHandler
}

func (s *Hub) Subscribe(handler MessageHandler) {
	s.subscribers = append(s.subscribers, handler)
}

func (s *Hub) Add(kind string, data any) {
	s.Start(kind, data)
	s.End()
}

func (s *Hub) Start(kind string, data any) {
	event := Message{Kind: kind, Data: data}
	event.Depth = s.stack.Size()
	s.events = append(s.events, event)
	s.stack.Push(event)
	s.Emit(event)
}

func (s *Hub) End() {
	if s.stack.Size() == 0 {
		return
	}
	last := s.stack.Pop()
	event := Message{Kind: last.Kind, Data: last.Data}
	event.Depth = s.stack.Size() + 1
	event.End = true
	s.events = append(s.events, event)
	s.Emit(event)
}

func (s *Hub) Emit(event Message) {
	for _, subscriber := range s.subscribers {
		subscriber(event)
	}
}
