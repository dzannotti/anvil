package eventemitter

import (
	"sync"
	"testing"
)

type TestEvent struct {
	Message string
}

func TestEventEmitter_New(t *testing.T) {
	emitter := New()
	if emitter == nil {
		t.Fatal("New() returned nil")
	}
	if emitter.handlers == nil {
		t.Fatal("handlers map not initialized")
	}
	if emitter.capturers == nil {
		t.Fatal("capturers slice not initialized")
	}
}

func TestEventEmitter_On(t *testing.T) {
	emitter := New()
	event := TestEvent{Message: "test"}
	called := false

	emitter.On(event, func(e any) {
		called = true
		if msg := e.(TestEvent).Message; msg != "test" {
			t.Errorf("expected message 'test', got '%s'", msg)
		}
	})

	emitter.Emit(event)
	if !called {
		t.Error("handler was not called")
	}
}

func TestEventEmitter_AddCapturer(t *testing.T) {
	emitter := New()
	event := TestEvent{Message: "test"}
	called := false

	emitter.AddCapturer(func(e any) {
		called = true
		if msg := e.(TestEvent).Message; msg != "test" {
			t.Errorf("expected message 'test', got '%s'", msg)
		}
	})

	emitter.Emit(event)
	if !called {
		t.Error("capturer was not called")
	}
}

func TestEventEmitter_MultipleHandlers(t *testing.T) {
	emitter := New()
	event := TestEvent{Message: "test"}
	callCount := 0

	handler := func(_ any) { callCount++ }
	emitter.On(event, handler)
	emitter.On(event, handler)
	emitter.Emit(event)

	if callCount != 2 {
		t.Errorf("expected 2 handler calls, got %d", callCount)
	}
}

func TestEventEmitter_ConcurrentEmit(t *testing.T) {
	emitter := New()
	event := TestEvent{Message: "test"}
	var wg sync.WaitGroup
	var counter int
	var mu sync.Mutex

	emitter.On(event, func(_ any) {
		mu.Lock()
		counter++
		mu.Unlock()
	})

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			emitter.Emit(event)
		}()
	}

	wg.Wait()
	if counter != 100 {
		t.Errorf("expected 100 emissions, got %d", counter)
	}
}

func TestEventEmitter_DifferentEventTypes(t *testing.T) {
	emitter := New()
	event1 := TestEvent{Message: "test1"}
	event2 := struct{ Value int }{Value: 42}
	callCount1 := 0
	callCount2 := 0

	emitter.On(event1, func(_ any) { callCount1++ })
	emitter.On(event2, func(_ any) { callCount2++ })

	emitter.Emit(event1)
	emitter.Emit(event2)

	if callCount1 != 1 || callCount2 != 1 {
		t.Errorf("expected one call for each event type, got %d and %d", callCount1, callCount2)
	}
}
