package effect

import (
	"anvil/internal/effect/state"
	"anvil/internal/expression"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	t.Parallel()

	t.Run("creates empty container", func(t *testing.T) {
		c := NewContainer()
		assert.NotNil(t, c)
		assert.Empty(t, c.effects)
	})

	t.Run("creates container with initial effects", func(t *testing.T) {
		effect1 := New("test1")
		effect2 := New("test2")
		c := NewContainer(effect1, effect2)
		assert.NotNil(t, c)
		assert.Len(t, c.effects, 2)
		assert.Contains(t, c.effects, effect1)
		assert.Contains(t, c.effects, effect2)
	})
}

func TestContainer_Add(t *testing.T) {
	t.Parallel()

	t.Run("adds effect to empty container", func(t *testing.T) {
		c := NewContainer()
		effect := New("test")
		c.Add(effect)
		assert.Len(t, c.effects, 1)
		assert.Contains(t, c.effects, effect)
	})

	t.Run("adds multiple effects", func(t *testing.T) {
		c := NewContainer()
		effect1 := New("test1")
		effect2 := New("test2")
		c.Add(effect1)
		c.Add(effect2)
		assert.Len(t, c.effects, 2)
		assert.Contains(t, c.effects, effect1)
		assert.Contains(t, c.effects, effect2)
	})
}

func TestContainer_Remove(t *testing.T) {
	t.Parallel()

	t.Run("removes existing effect", func(t *testing.T) {
		effect1 := New("test1")
		effect2 := New("test2")
		c := NewContainer(effect1, effect2)
		c.Remove(effect1)
		assert.Len(t, c.effects, 1)
		assert.NotContains(t, c.effects, effect1)
		assert.Contains(t, c.effects, effect2)
	})

	t.Run("does nothing when removing non-existent effect", func(t *testing.T) {
		effect1 := New("test1")
		effect2 := New("test2")
		c := NewContainer(effect1)
		c.Remove(effect2)
		assert.Len(t, c.effects, 1)
		assert.Contains(t, c.effects, effect1)
	})
}

func TestContainer_Evaluate(t *testing.T) {
	t.Run("evaluates all effects", func(t *testing.T) {
		evaluationCount := 0
		effect1 := New("test1", WithAttributeCalculation(func(e *Effect, state *state.AttributeCalculation, wg *sync.WaitGroup) {
			evaluationCount++
			state.Expression.AddScalar(2, "test")
			wg.Done()
		}))
		effect2 := New("test2", WithAttributeCalculation(func(e *Effect, state *state.AttributeCalculation, wg *sync.WaitGroup) {
			evaluationCount++
			wg.Done()
		}))

		c := NewContainer(effect1, effect2)
		state := &state.AttributeCalculation{
			Expression: expression.FromScalar(5, "bar"),
		}
		wg := &sync.WaitGroup{}
		wg.Add(1)
		c.Evaluate(state, wg)
		wg.Wait()
		res := state.Expression.Evaluate()
		assert.Equal(t, 2, evaluationCount)
		assert.Equal(t, 7, res.Value)
	})

	t.Run("handles empty container", func(t *testing.T) {
		c := NewContainer()
		state := &state.AttributeCalculation{}
		wg := &sync.WaitGroup{}
		wg.Add(1)
		c.Evaluate(state, wg) // Should not panic
		wg.Wait()
		assert.Equal(t, true, true)
	})
}
