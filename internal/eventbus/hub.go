package eventbus

import (
	"anvil/internal/collection"
)

type MessageHandler func(data Message)

type Hub struct {
	events      []Message
	stack       collection.Stack[Message]
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
	event.Depth = s.stack.Len()
	s.events = append(s.events, event)
	s.stack.Push(event)
	s.Emit(event)
}

func (s *Hub) End() {
	last, ok := s.stack.Pop()
	if !ok {
		return
	}
	event := Message{Data: last.Data}
	event.Depth = s.stack.Len() + 1
	event.IsEnd = true
	s.events = append(s.events, event)
	s.Emit(event)
}

func (s *Hub) Emit(event Message) {
	for _, subscriber := range s.subscribers {
		subscriber(event)
	}
}
