package effect

import (
	"sync"
)

type Priority int

const (
	PriorityNormal       Priority = iota
	PriorityEarly        Priority = -20
	PriorityBase         Priority = -60
	PriorityBaseOverride Priority = -40
	PriorityLate         Priority = 20
	PriorityLast         Priority = 40
)

type Effect struct {
	Name     string
	Handlers map[string]func(*Effect, any, *sync.WaitGroup)
	Priority Priority
}

func (e *Effect) Evaluate(event string, state any) {
	handler, exists := e.Handlers[event]
	if exists {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go handler(e, state, wg)
		wg.Wait()
	}
}

func (e *Effect) WithHandler(event string, handler func(*Effect, any, *sync.WaitGroup)) {
	e.Handlers[event] = func(e *Effect, s any, wg *sync.WaitGroup) {
		handler(e, s, wg)
	}
}
