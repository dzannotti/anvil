package effect

import (
	"anvil/internal/expression"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testState struct {
	Expression *expression.Expression
}

func TestEffect_Evaluate(t *testing.T) {
	t.Run("handler executes when event matches", func(t *testing.T) {
		e := &Effect{}
		called := false
		e.WithHandler("test", func(_ *Effect, _ any) {
			called = true
		})
		e.Evaluate("test", nil)
		assert.True(t, called)
	})

	t.Run("handler does not execute when event does not match", func(t *testing.T) {
		e := &Effect{}
		called := false
		e.WithHandler("test", func(_ *Effect, _ any) {
			called = true
		})

		e.Evaluate("other", nil)
		assert.False(t, called)
	})

	t.Run("state is modified by handler", func(t *testing.T) {
		e := &Effect{}
		expr := expression.FromScalar(10, "test")
		state := &testState{Expression: &expr}
		e.WithHandler("modify", func(_ *Effect, s any) {
			state := s.(*testState) // panics if not the correct typ
			state.Expression.AddScalar(5, "test")
		})
		e.Evaluate("modify", state)
		state.Expression.Evaluate()
		assert.Equal(t, 15, state.Expression.Value)
	})
}
