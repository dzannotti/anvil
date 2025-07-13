package expression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/expression"
	"anvil/internal/tag"
)

func TestDiceComponent_Evaluate(t *testing.T) {
	t.Run("rolls dice and sums results", func(t *testing.T) {
		expr := expression.FromDice(3, 6, "test")
		expr.Rng = newMockRoller(2, 4, 6)
		expr.Evaluate()

		assert.Equal(t, 12, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, expression.ComponentKindDice, comp.Kind())
		assert.Equal(t, 12, comp.Value())
		assert.Equal(t, "test", comp.Source())
	})

	t.Run("handles single die", func(t *testing.T) {
		expr := expression.FromDice(1, 20, "single")
		expr.Rng = newMockRoller(15)
		expr.Evaluate()

		assert.Equal(t, 15, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, 15, comp.Value())
	})

	t.Run("handles negative dice count", func(t *testing.T) {
		expr := expression.FromDice(-2, 6, "penalty")
		expr.Rng = newMockRoller(3, 4)
		expr.Evaluate()

		assert.Equal(t, -7, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, -7, comp.Value())
	})

	t.Run("handles zero dice", func(t *testing.T) {
		expr := expression.FromDice(0, 6, "none")
		expr.Evaluate()

		assert.Equal(t, 0, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, 0, comp.Value())
	})

	t.Run("damage dice has correct tags", func(t *testing.T) {
		tags := tag.ContainerFromString("damage.fire", "damage.spell")
		expr := expression.FromDamageDice(2, 8, tags, "flame strike")
		expr.Rng = newMockRoller(5, 7)
		expr.Evaluate()

		assert.Equal(t, 12, expr.Value)
		comp := expr.Components[0]
		compTags := comp.Tags()
		assert.True(t, compTags.HasTag(tag.FromString("damage.fire")))
		assert.True(t, compTags.HasTag(tag.FromString("damage.spell")))
	})

	t.Run("component provides dice info through interface", func(t *testing.T) {
		expr := expression.FromDice(3, 8, "weapon")
		comp := expr.Components[0]

		diceComp, ok := comp.(*expression.DiceComponent)
		assert.True(t, ok)
		assert.Equal(t, 3, diceComp.Times())
		assert.Equal(t, 8, diceComp.Sides())
	})
}
