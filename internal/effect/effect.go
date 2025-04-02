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

type Handlers map[string]func(*Effect, any)

func (h *Handlers) get() Handlers {
	if *h == nil {
		*h = make(map[string]func(*Effect, any))
	}
	return *h
}

type Effect struct {
	Name     string
	Handlers Handlers
	Priority Priority
}

func (e *Effect) Evaluate(event string, state any) {
	handler, exists := e.Handlers.get()[event]
	if exists {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler(e, state)
		}()
		wg.Wait()
	}
}

func (e *Effect) WithHandler(event string, handler func(*Effect, any)) {
	e.Handlers.get()[event] = handler
}
