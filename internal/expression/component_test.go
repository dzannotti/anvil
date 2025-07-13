package expression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/expression"
	"anvil/internal/tag"
)

func TestComponent_Interface(t *testing.T) {
	t.Run("constant component implements interface", func(t *testing.T) {
		expr := expression.FromConstant(5, "test")
		comp := expr.Components[0]

		assert.Equal(t, expression.ComponentKindConstant, comp.Kind())
		assert.Equal(t, 5, comp.Value())
		assert.Equal(t, "test", comp.Source())
		tags := comp.Tags()
		assert.True(t, tags.HasTag(tag.FromString("primary")))
		assert.Empty(t, comp.Components())
	})

	t.Run("dice component implements interface", func(t *testing.T) {
		expr := expression.FromDice(2, 6, "test")
		comp := expr.Components[0]

		assert.Equal(t, expression.ComponentKindDice, comp.Kind())
		assert.Equal(t, 0, comp.Value())
		assert.Equal(t, "test", comp.Source())
		tags := comp.Tags()
		assert.True(t, tags.HasTag(tag.FromString("primary")))
		assert.Empty(t, comp.Components())
	})

	t.Run("d20 component implements interface", func(t *testing.T) {
		expr := expression.FromD20("test")
		comp := expr.Components[0]

		assert.Equal(t, expression.ComponentKindD20, comp.Kind())
		assert.Equal(t, 0, comp.Value())
		assert.Equal(t, "test", comp.Source())
		tags := comp.Tags()
		assert.True(t, tags.HasTag(tag.FromString("primary")))
		assert.Empty(t, comp.Components())
	})
}
