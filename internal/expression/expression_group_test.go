package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tagcontainer"
)

func TestExpression_ExpressionGroup(t *testing.T) {
	t.Run("can evaluate a scalar", func(t *testing.T) {
		expression := Expression{}
		expression.AddDamageScalar(2, "Damage", *tagcontainer.FromString("Slashing"))
		res := expression.EvaluateGroup()
		assert.Equal(t, 2, res.Value)
		assert.Equal(t, len(res.Terms), 1)
		assert.Equal(t, res.Terms[0].Tags, tagcontainer.FromString("Slashing"))
	})
	t.Run("groups by tags", func(t *testing.T) {
		expression := Expression{}
		expression.AddDamageScalar(2, "Damage", *tagcontainer.FromString("Slashing"))
		expression.AddDamageScalar(4, "Damage", *tagcontainer.FromString("Slashing"))
		expression.AddDamageScalar(3, "Damage", *tagcontainer.FromString("Magical"))
		res := expression.EvaluateGroup()
		assert.Equal(t, 9, res.Value)
		assert.Equal(t, len(res.Terms), 2)
		assert.Equal(t, res.Terms[0].Tags, tagcontainer.FromString("Slashing"))
		assert.Equal(t, res.Terms[1].Tags, tagcontainer.FromString("Magical"))
	})

	t.Run("can evaluate a dice", func(t *testing.T) {
		expression := Expression{rng: &mockRoller{mockReturns: []int{5}}}
		expression.AddDamageDice(1, 6, "Damage", *tagcontainer.FromString("Slashing"))
		expression.AddDamageScalar(2, "Damage", *tagcontainer.FromString("Slashing"))
		res := expression.EvaluateGroup()
		assert.Equal(t, 7, res.Value)
		assert.Equal(t, len(res.Terms), 1)
		assert.Equal(t, res.Terms[0].Tags, tagcontainer.FromString("Slashing"))
	})

	t.Run("should halve damage for match tags", func(t *testing.T) {
		expression := Expression{rng: &mockRoller{mockReturns: []int{5}}}
		expression.AddDamageDice(1, 6, "Damage", *tagcontainer.FromStrings([]string{"Slashing", "fire"}))
		expression.AddDamageScalar(2, "Damage", *tagcontainer.FromStrings([]string{"Slashing", "fire"}))
		expression.ReplaceWith(3, "Bad Luck")
		res := expression.EvaluateGroup()
		assert.Len(t, expression.Terms, 1)
		assert.Equal(t, 3, res.Value)
	})

	t.Run("should double dice damage", func(t *testing.T) {
		expression := Expression{rng: &mockRoller{mockReturns: []int{5, 3, 20, 18}}}
		expression.AddDamageDice(1, 6, "Damage", *tagcontainer.FromStrings([]string{"Slashing", "fire"}))
		expression.AddDamageScalar(2, "Damage", *tagcontainer.FromStrings([]string{"Slashing", "fire"}))
		expression.AddDamageDice(1, 6, "Damage", *tagcontainer.FromString("Slashing"))
		expression.AddDamageScalar(4, "Damage", *tagcontainer.FromString("slashing"))
		expression.DoubleDice("Critical")
		res := expression.EvaluateGroup()
		assert.Equal(t, 52, res.Value)
	})

	t.Run("should max dice damage", func(t *testing.T) {
		expression := Expression{rng: &mockRoller{mockReturns: []int{5, 3}}}
		expression.AddDamageDice(1, 6, "Damage", *tagcontainer.FromStrings([]string{"Slashing", "fire"}))
		expression.AddDamageScalar(2, "Damage", *tagcontainer.FromStrings([]string{"Slashing", "fire"}))
		expression.AddDamageDice(1, 6, "Damage", *tagcontainer.FromString("Slashing"))
		expression.AddDamageScalar(4, "Damage", *tagcontainer.FromString("slashing"))
		expression.MaxDice("Critical")
		res := expression.EvaluateGroup()
		assert.Equal(t, 26, res.Value)
	})
}
