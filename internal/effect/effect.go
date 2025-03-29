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

func (e *Effect) Evaluate(state state.State) {
	handler, exists := e.handlers[state.Type()]
	if exists {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go handler(e, state, wg)
		wg.Wait()
	}
}

func withHandler[S state.State](stateType state.Type, handler func(*Effect, S, *sync.WaitGroup)) Option {
	return func(e Effect) {
		e.handlers[stateType] = func(e *Effect, s state.State, wg *sync.WaitGroup) {
			if state, ok := s.(S); ok {
				handler(e, state, wg)
			}
		}
	}
}

func WithAttributeCalculation(handler func(*Effect, *state.AttributeCalculation, *sync.WaitGroup)) Option {
	return withHandler(state.AttributeCalculationType, handler)
}

func WithBeforeAttackRoll(handler func(*Effect, *state.BeforeAttackRoll, *sync.WaitGroup)) Option {
	return withHandler(state.BeforeAttackRollType, handler)
}

func WithAfterAttackRoll(handler func(*Effect, *state.AfterAttackRoll, *sync.WaitGroup)) Option {
	return withHandler(state.AfterAttackRollType, handler)
}
