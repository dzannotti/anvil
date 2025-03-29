package effect

import (
	"anvil/internal/effect/state"
	"anvil/internal/expression"
	"anvil/internal/tag"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEffect_NewEffect(t *testing.T) {
	t.Parallel()
	t.Run("can create an Effect", func(t *testing.T) {
		effect := New("test")
		assert.NotNil(t, effect)
		assert.NotNil(t, effect.handlers)
		assert.Empty(t, effect.handlers)
	})

	t.Run("can create an attribute calculation effect", func(t *testing.T) {
		effect := New("test", WithAttributeCalculation(func(e *Effect, state *state.AttributeCalculation, wg *sync.WaitGroup) { wg.Done() }))
		assert.NotNil(t, effect)
		assert.NotNil(t, effect.handlers)
		assert.Len(t, effect.handlers, 1)
		assert.Contains(t, effect.handlers, state.AttributeCalculationType)
	})
}

func TestEffect_Evaluate(t *testing.T) {
	t.Run("can evaluate with no handlers", func(t *testing.T) {
		effect := New("with")
		state := &state.AttributeCalculation{}
		effect.Evaluate(state) // Should not panic
		assert.Equal(t, true, true)
	})

	t.Run("can evaluate with matching handler", func(t *testing.T) {
		handlerCalled := false
		effect := New("test", WithAttributeCalculation(func(e *Effect, state *state.AttributeCalculation, wg *sync.WaitGroup) {
			handlerCalled = true
			state.Expression.AddScalar(2, "test")
			state.Attribute = tag.FromString("foo")
			wg.Done()
		}))
		state := &state.AttributeCalculation{
			Expression: expression.FromScalar(5, "bar"),
		}
		effect.Evaluate(state)
		assert.True(t, handlerCalled)
		assert.Equal(t, state.Attribute, tag.FromString("foo"))
		assert.Equal(t, 7, state.Expression.Evaluate().Value, tag.FromString("foo"))
	})

	t.Run("does not evaluate with non-matching handler", func(t *testing.T) {
		handlerCalled := false
		effect := New("test", WithAttributeCalculation(func(e *Effect, state *state.AttributeCalculation, wg *sync.WaitGroup) {
			handlerCalled = true
			wg.Done()
		}))
		state := &state.BeforeAttackRoll{}
		effect.Evaluate(state)
		assert.False(t, handlerCalled)
	})
}
