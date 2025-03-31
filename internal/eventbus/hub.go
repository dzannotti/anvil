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

func (s *Hub) Add(data any) {
	s.Start(data)
	s.End()
}

func (s *Hub) Start(data any) {
	event := Message{Data: data}
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
	event := Message{Data: last.Data}
	event.Depth = s.stack.Size() + 1
	event.IsEnd = true
	s.events = append(s.events, event)
	s.Emit(event)
}

func (s *Hub) Emit(event Message) {
	for _, subscriber := range s.subscribers {
		subscriber(event)
	}
}
