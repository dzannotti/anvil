package core

import (
	"testing"

	"anvil/internal/expression"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testState struct {
	Expression *expression.Expression
}

type TestEffect struct {
	Expression *expression.Expression
}

func TestEffect_Evaluate(t *testing.T) {
	t.Run("handler executes when event matches", func(t *testing.T) {
		e := &Effect{}
		called := false
		e.withHandler("TestEffect", func(_ *Effect, _ any) {
			called = true
		})
		e.Evaluate(&TestEffect{})
		assert.True(t, called)
	})

	t.Run("handler does not execute when event does not match", func(t *testing.T) {
		e := &Effect{}
		called := false
		e.withHandler("TestEffect", func(_ *Effect, _ any) {
			called = true
		})

		e.Evaluate(&testState{})
		assert.False(t, called)
	})

	t.Run("state is modified by handler", func(t *testing.T) {
		e := &Effect{}
		expr := expression.FromConstant(10, "test")
		state := &TestEffect{Expression: &expr}
		e.withHandler("TestEffect", func(_ *Effect, s any) {
			estate, ok := s.(*TestEffect)
			require.True(t, ok, "state must be *TestEffect for test to continue")
			estate.Expression.AddConstant(5, "test")
		})
		e.Evaluate(state)
		state.Expression.Evaluate()
		assert.Equal(t, 15, state.Expression.Value)
	})
}
