package expression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/expression"
	"anvil/internal/tag"
)

func TestConstantComponent_Evaluate(t *testing.T) {
	t.Run("returns constant value", func(t *testing.T) {
		expr := expression.FromConstant(42, "test")
		expr.Evaluate()

		assert.Equal(t, 42, expr.Value)
		assert.Len(t, expr.Components, 1)

		comp := expr.Components[0]
		assert.Equal(t, expression.ComponentKindConstant, comp.Kind())
		assert.Equal(t, 42, comp.Value())
		assert.Equal(t, "test", comp.Source())
		tags := comp.Tags()
		assert.True(t, tags.HasTag(tag.FromString("primary")))
	})

	t.Run("works with negative values", func(t *testing.T) {
		expr := expression.FromConstant(-5, "penalty")
		expr.Evaluate()

		assert.Equal(t, -5, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, -5, comp.Value())
	})

	t.Run("works with zero", func(t *testing.T) {
		expr := expression.FromConstant(0, "zero")
		expr.Evaluate()

		assert.Equal(t, 0, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, 0, comp.Value())
	})

	t.Run("damage constant has correct tags", func(t *testing.T) {
		tags := tag.ContainerFromString("damage.fire", "damage.elemental")
		expr := expression.FromDamageConstant(10, tags, "fireball")
		expr.Evaluate()

		assert.Equal(t, 10, expr.Value)
		comp := expr.Components[0]
		compTags := comp.Tags()
		assert.True(t, compTags.HasTag(tag.FromString("damage.fire")))
		assert.True(t, compTags.HasTag(tag.FromString("damage.elemental")))
	})
}
