package log

import (
	"anvil/internal/collection"
	"anvil/internal/eventemitter"
)

type EventHandler func(data Event)

type EventLog struct {
	events  []Event
	stack   collection.Stack[Event]
	emitter *eventemitter.EventEmitter
}

func New() *EventLog {
	return &EventLog{
		emitter: eventemitter.New(),
	}
}

func (l *EventLog) AddCapturer(handler EventHandler) {
	l.emitter.AddCapturer(func(event any) {
		handler(event.(Event))
	})
}

func (l *EventLog) Add(data any) {
	l.Start(data)
	l.End()
}

func (l *EventLog) Start(data any) {
	event := NewEvent(data)
	event.Depth = l.stack.Len()
	l.events = append(l.events, event)
	l.stack.Push(event)
	l.emitter.Emit(event)
}

func (l *EventLog) End() {
	last, ok := l.stack.Pop()
	if !ok {
		return
	}
	event := NewEvent(last.Data)
	event.Depth = l.stack.Len() + 1
	event.IsEnd = true
	l.events = append(l.events, event)
	l.emitter.Emit(event)
}
