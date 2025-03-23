package collection

import "testing"

func TestStack(t *testing.T) {
	t.Run("new stack should be empty", func(t *testing.T) {
		s := New[int](10)
		if s.Len() != 0 {
			t.Errorf("expected empty stack, got length %d", s.Len())
		}
	})

	t.Run("push should increase length", func(t *testing.T) {
		s := New[int](10)
		s.Push(1)
		if s.Len() != 1 {
			t.Errorf("expected length 1, got %d", s.Len())
		}
	})

	t.Run("pop from empty stack returns false", func(t *testing.T) {
		s := New[int](10)
		_, ok := s.Pop()
		if ok {
			t.Error("expected false when popping from empty stack")
		}
	})

	t.Run("push and pop should maintain LIFO order", func(t *testing.T) {
		s := New[int](10)
		s.Push(1)
		s.Push(2)
		s.Push(3)

		expected := []int{3, 2, 1}
		for i, want := range expected {
			got, ok := s.Pop()
			if !ok {
				t.Errorf("pop %d: expected success, got failure", i)
			}
			if got != want {
				t.Errorf("pop %d: expected %d, got %d", i, want, got)
			}
		}
	})

	t.Run("length should decrease after pop", func(t *testing.T) {
		s := New[int](10)
		s.Push(1)
		s.Push(2)
		initialLen := s.Len()

		s.Pop()
		if s.Len() != initialLen-1 {
			t.Errorf("expected length %d, got %d", initialLen-1, s.Len())
		}
	})

	t.Run("works with string type", func(t *testing.T) {
		s := New[string](10)
		s.Push("hello")
		s.Push("world")

		got, ok := s.Pop()
		if !ok {
			t.Error("expected successful pop")
		}
		if got != "world" {
			t.Errorf("expected 'world', got '%s'", got)
		}
	})

	t.Run("respects initial capacity", func(t *testing.T) {
		s := New[int](2)
		s.Push(1)
		s.Push(2)
		s.Push(3) // Should still work, just reallocates

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})
}
