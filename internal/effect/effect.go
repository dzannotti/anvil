package effect

import (
	"anvil/internal/effect/state"
	"sync"
)

type EffectOption func(Effect)

type Effect struct {
	Id       string
	handlers map[state.Type]func(*Effect, state.State, *sync.WaitGroup)
}

func New(id string, opts ...EffectOption) Effect {
	e := Effect{
		Id:       id,
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
		wg.Add(1)
		go handler(e, state, wg)
	}
}

func WithAttributeCalculation(handler func(*Effect, *state.AttributeCalculation, *sync.WaitGroup)) EffectOption {
	return func(e Effect) {
		e.handlers[state.AttributeCalculationType] = func(e *Effect, s state.State, wg *sync.WaitGroup) {
			if attributeState, ok := s.(*state.AttributeCalculation); ok {
				handler(e, attributeState, wg)
			}
		}
	}
}
