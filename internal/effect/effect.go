package effect

import (
	"anvil/internal/effect/state"
	"sync"
)

type Option func(Effect)

type Effect struct {
	Name     string
	handlers map[state.Type]func(*Effect, state.State, *sync.WaitGroup)
}

func New(name string, opts ...Option) Effect {
	e := Effect{
		Name:     name,
		handlers: make(map[state.Type]func(*Effect, state.State, *sync.WaitGroup)),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Effect) Evaluate(state state.State, wg *sync.WaitGroup) {
	handler, exists := e.handlers[state.Type()]
	if exists {
		ewg := &sync.WaitGroup{}
		ewg.Add(1)
		go handler(e, state, ewg)
		ewg.Wait()
	}
	wg.Done()
}

func WithAttributeCalculation(handler func(*Effect, *state.AttributeCalculation, *sync.WaitGroup)) Option {
	return func(e Effect) {
		e.handlers[state.AttributeCalculationType] = func(e *Effect, s state.State, wg *sync.WaitGroup) {
			if attributeState, ok := s.(*state.AttributeCalculation); ok {
				handler(e, attributeState, wg)
			}
		}
	}
}
