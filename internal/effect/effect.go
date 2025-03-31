package effect

import (
	"sync"
)

type EffectPriority int

const (
	PriorityNormal       EffectPriority = iota
	PriorityEarly        EffectPriority = -20
	PriorityBase         EffectPriority = -60
	PriorityBaseOverride EffectPriority = -40
	PriorityLate         EffectPriority = 20
	PriorityLast         EffectPriority = 40
)

type Effect struct {
	Name     string
	Handlers map[string]func(*Effect, any, *sync.WaitGroup)
	Priority EffectPriority
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
