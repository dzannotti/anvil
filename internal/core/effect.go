package core

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

func (e *Effect) withHandler(event string, handler func(*Effect, any)) {
	e.Handlers.get()[event] = handler
}

func (e *Effect) WithBeforeAttackRoll(handler func(*Effect, *BeforeAttackRollState)) {
	e.Handlers.get()[BeforeAttackRoll] = func(e *Effect, state any) {
		handler(e, state.(*BeforeAttackRollState))
	}
}

func (e *Effect) WithAfterAttackRoll(handler func(*Effect, *AfterAttackRollState)) {
	e.Handlers.get()[AfterAttackRoll] = func(e *Effect, state any) {
		handler(e, state.(*AfterAttackRollState))
	}
}

func (e *Effect) WithAttributeCalculation(handler func(*Effect, *AttributeCalculationState)) {
	e.Handlers.get()[AttributeCalculation] = func(e *Effect, state any) {
		handler(e, state.(*AttributeCalculationState))
	}
}

func (e *Effect) WithBeforeTakeDamage(handler func(*Effect, *BeforeTakeDamageState)) {
	e.Handlers.get()[BeforeTakeDamage] = func(e *Effect, state any) {
		handler(e, state.(*BeforeTakeDamageState))
	}
}

func (e *Effect) WithAfterTakeDamage(handler func(*Effect, *AfterTakeDamageState)) {
	e.Handlers.get()[AfterTakeDamage] = func(e *Effect, state any) {
		handler(e, state.(*AfterTakeDamageState))
	}
}

func (e *Effect) WithBeforeDamageRoll(handler func(*Effect, *BeforeDamageRollState)) {
	e.Handlers.get()[BeforeDamageRoll] = func(e *Effect, state any) {
		handler(e, state.(*BeforeDamageRollState))
	}
}

func (e *Effect) WithAfterDamageRoll(handler func(*Effect, *AfterDamageRollState)) {
	e.Handlers.get()[AfterDamageRoll] = func(e *Effect, state any) {
		handler(e, state.(*AfterDamageRollState))
	}
}

func (e *Effect) WithBeforeSavingThrow(handler func(*Effect, *BeforeSavingThrowState)) {
	e.Handlers.get()[BeforeSavingThrow] = func(e *Effect, state any) {
		handler(e, state.(*BeforeSavingThrowState))
	}
}

func (e *Effect) WithAfterSavingThrow(handler func(*Effect, *AfterSavingThrowState)) {
	e.Handlers.get()[AfterSavingThrow] = func(e *Effect, state any) {
		handler(e, state.(*AfterSavingThrowState))
	}
}

func (e *Effect) WithAttributeChanged(handler func(*Effect, *AttributeChangedState)) {
	e.Handlers.get()[AttributeChanged] = func(e *Effect, state any) {
		handler(e, state.(*AttributeChangedState))
	}
}

func (e *Effect) WithTurnStarted(handler func(*Effect, *TurnState)) {
	e.Handlers.get()[TurnStarted] = func(e *Effect, state any) {
		handler(e, state.(*TurnState))
	}
}

func (e *Effect) WithTurnEnded(handler func(*Effect, *TurnState)) {
	e.Handlers.get()[TurnEnded] = func(e *Effect, state any) {
		handler(e, state.(*TurnState))
	}
}

func (e *Effect) WithConditionAdded(handler func(*Effect, *ConditionChangedState)) {
	e.Handlers.get()[ConditionAdded] = func(e *Effect, state any) {
		handler(e, state.(*ConditionChangedState))
	}
}

func (e *Effect) WithConditionRemoved(handler func(*Effect, *ConditionChangedState)) {
	e.Handlers.get()[ConditionRemoved] = func(e *Effect, state any) {
		handler(e, state.(*ConditionChangedState))
	}
}
