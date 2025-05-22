package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_ExpressionGroup(t *testing.T) {
	t.Run("can evaluate a scalar", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageScalar(2, "Damage", tag.ContainerFromString("Slashing"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 2, res.Value)
		assert.Len(t, res.Terms, 1)
		assert.Equal(t, res.Terms[0].Tags, tag.ContainerFromString("Slashing"))
	})
	t.Run("groups by tags", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageScalar(2, "Damage", tag.ContainerFromString("Slashing"))
		expr.AddDamageScalar(4, "Damage", tag.ContainerFromString("Slashing"))
		expr.AddDamageScalar(3, "Damage", tag.ContainerFromString("Magical"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 9, res.Value)
		assert.Len(t, res.Terms, 2)
		assert.Equal(t, res.Terms[0].Tags, tag.ContainerFromString("Slashing"))
		assert.Equal(t, res.Terms[1].Tags, tag.ContainerFromString("Magical"))
	})

	t.Run("can evaluate a dice", func(t *testing.T) {
		expr := Expression{rng: &mockRoller{mockReturns: []int{5}}}
		expr.AddDamageDice(1, 6, "Damage", tag.ContainerFromString("Slashing"))
		expr.AddDamageScalar(2, "Damage", tag.ContainerFromString("Slashing"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 7, res.Value)
		assert.Len(t, res.Terms, 1)
		assert.Equal(t, res.Terms[0].Tags, tag.ContainerFromString("Slashing"))
	})

	t.Run("should halve damage for match tags", func(t *testing.T) {
		expr := Expression{rng: &mockRoller{mockReturns: []int{5}}}
		expr.AddDamageDice(1, 6, "Damage", tag.ContainerFromStrings([]string{"Slashing", "fire"}))
		expr.AddDamageScalar(2, "Damage", tag.ContainerFromStrings([]string{"Slashing", "fire"}))
		expr.ReplaceWith(3, "Bad Luck")
		res := expr.EvaluateGroup()
		assert.Len(t, expr.Terms, 1)
		assert.Equal(t, 3, res.Value)
	})

	t.Run("should double dice damage", func(t *testing.T) {
		expr := Expression{rng: &mockRoller{mockReturns: []int{5, 3, 20, 18}}}
		expr.AddDamageDice(1, 6, "Damage", tag.ContainerFromStrings([]string{"Slashing", "fire"}))
		expr.AddDamageScalar(2, "Damage", tag.ContainerFromStrings([]string{"Slashing", "fire"}))
		expr.AddDamageDice(1, 6, "Damage", tag.ContainerFromString("Slashing"))
		expr.AddDamageScalar(4, "Damage", tag.ContainerFromString("slashing"))
		expr.DoubleDice("Critical")
		res := expr.EvaluateGroup()
		assert.Equal(t, 52, res.Value)
	})

	t.Run("should max dice damage", func(t *testing.T) {
		expr := Expression{rng: &mockRoller{mockReturns: []int{5, 3}}}
		expr.AddDamageDice(1, 6, "Damage", tag.ContainerFromStrings([]string{"Slashing", "fire"}))
		expr.AddDamageScalar(2, "Damage", tag.ContainerFromStrings([]string{"Slashing", "fire"}))
		expr.AddDamageDice(1, 6, "Damage", tag.ContainerFromString("Slashing"))
		expr.AddDamageScalar(4, "Damage", tag.ContainerFromString("slashing"))
		expr.MaxDice("Critical")
		res := expr.EvaluateGroup()
		assert.Equal(t, 26, res.Value)
	})
}
