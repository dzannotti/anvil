package core

import (
	"reflect"
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

func (e *Effect) Evaluate(state any) {
	stateType := reflect.TypeOf(state)
	if stateType.Kind() != reflect.Ptr {
		panic("state must be a pointer")
	}

	// Get the event name from the struct type (remove pointer)
	eventName := stateType.Elem().Name()

	handler, exists := e.Handlers.get()[eventName]
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

// On registers a handler using reflection to derive the event name from the state type
func (e *Effect) On(handler any) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	if handlerType.NumIn() != 1 {
		panic("handler must have exactly one parameter")
	}

	paramType := handlerType.In(0)
	if paramType.Kind() != reflect.Ptr {
		panic("handler parameter must be a pointer")
	}

	// Get the event name from the struct type (remove pointer and "State" suffix)
	eventName := paramType.Elem().Name()

	// Convert handler to the internal signature
	handlerValue := reflect.ValueOf(handler)
	e.Handlers.get()[eventName] = func(_ *Effect, state any) {
		handlerValue.Call([]reflect.Value{reflect.ValueOf(state)})
	}
}
