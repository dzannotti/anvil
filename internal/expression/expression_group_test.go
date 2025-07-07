package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_ExpressionGroup(t *testing.T) {
	t.Run("can evaluate a constant", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 2, res.Value)
		assert.Len(t, res.Components, 1)
		assert.Equal(t, res.Components[0].Tags, tag.NewContainerFromString("Slashing"))
	})
	t.Run("groups by tags", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(4, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(3, "Damage", tag.NewContainerFromString("Magical"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 9, res.Value)
		assert.Len(t, res.Components, 2)
		assert.Equal(t, res.Components[0].Tags, tag.NewContainerFromString("Slashing"))
		assert.Equal(t, res.Components[1].Tags, tag.NewContainerFromString("Magical"))
	})

	t.Run("can evaluate a dice", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{5}}}
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 7, res.Value)
		assert.Len(t, res.Components, 1)
		assert.Equal(t, res.Components[0].Tags, tag.NewContainerFromString("Slashing"))
	})

	t.Run("should halve damage for match tags", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{5}}}
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing", "fire"))
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing", "fire"))
		expr.ReplaceWith(3, "Bad Luck")
		res := expr.EvaluateGroup()
		assert.Len(t, expr.Components, 1)
		assert.Equal(t, 3, res.Value)
	})

	t.Run("should double dice damage", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{5, 3, 20, 18}}}
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing", "fire"))
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing", "fire"))
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(4, "Damage", tag.NewContainerFromString("slashing"))
		expr.DoubleDice("Critical")
		res := expr.EvaluateGroup()
		assert.Equal(t, 52, res.Value)
	})

	t.Run("should max dice damage", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{5, 3}}}
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing", "fire"))
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing", "fire"))
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(4, "Damage", tag.NewContainerFromString("slashing"))
		expr.MaxDice("Critical")
		res := expr.EvaluateGroup()
		assert.Equal(t, 26, res.Value)
	})
}
