package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_PrimaryDamageCollation(t *testing.T) {
	t.Run("flaming sword with empty tag inheritance", func(t *testing.T) {
		// Flaming longsword: 1d8 slashing + 2d4 fire + STR modifier
		expr := Expression{Rng: &mockRoller{mockReturns: []int{8, 4, 2}}}             // d8=8, d4=4,2
		expr.AddDamageDice(1, 8, "Longsword", tag.NewContainerFromString("slashing")) // Primary damage
		expr.AddDamageDice(2, 4, "Flame Tongue", tag.NewContainerFromString("fire"))  // Explicit fire
		expr.AddDamageConstant(3, "STR Modifier", tag.NewContainer())                 // Empty tags - should inherit slashing

		grouped := expr.EvaluateGroup()

		// Should result in exactly 2 root components grouped by damage type
		assert.Len(t, grouped.Components, 2, "should have 2 damage type groups")
		assert.Equal(t, 17, grouped.Value, "total damage should be 8+6+3 = 17")

		// First component should be slashing damage (longsword + STR modifier)
		slashingComp := grouped.Components[0]
		assert.Equal(t, 11, slashingComp.Value, "slashing damage should be 8+3 = 11")
		assert.True(t, slashingComp.Tags.HasAny(tag.NewContainerFromString("slashing")), "first component should have slashing tags")

		// Second component should be fire damage (flame tongue)
		fireComp := grouped.Components[1]
		assert.Equal(t, 6, fireComp.Value, "fire damage should be 4+2 = 6")
		assert.True(t, fireComp.Tags.HasAny(tag.NewContainerFromString("fire")), "second component should have fire tags")
	})

	t.Run("primary tag explicit inheritance", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{6, 3}}}
		expr.AddDamageDice(1, 8, "Weapon", tag.NewContainerFromString("slashing"))
		expr.AddDamageDice(1, 6, "Cold", tag.NewContainerFromString("cold"))
		expr.AddDamageConstant(3, "STR Mod", tag.NewContainerFromString("primary")) // Explicit primary tag

		grouped := expr.EvaluateGroup()

		assert.Len(t, grouped.Components, 2, "should have 2 damage type groups")

		// STR modifier should be grouped with slashing (first component)
		slashingComp := grouped.Components[0]
		assert.Equal(t, 9, slashingComp.Value, "slashing should be 6+3 = 9")
		assert.True(t, slashingComp.Tags.HasAny(tag.NewContainerFromString("slashing")))

		coldComp := grouped.Components[1]
		assert.Equal(t, 3, coldComp.Value, "cold should be 3")
		assert.True(t, coldComp.Tags.HasAny(tag.NewContainerFromString("cold")))
	})

	t.Run("AddConstant also inherits from primary damage", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{5}}}
		expr.AddDamageDice(1, 6, "Weapon", tag.NewContainerFromString("piercing"))
		expr.AddConstant(2, "DEX Modifier") // AddConstant with empty tags should also inherit

		grouped := expr.EvaluateGroup()

		assert.Len(t, grouped.Components, 1, "should have 1 damage type group")
		assert.Equal(t, 7, grouped.Components[0].Value, "piercing damage should be 5+2 = 7")
		assert.True(t, grouped.Components[0].Tags.HasAny(tag.NewContainerFromString("piercing")))
	})

	t.Run("mixed AddDamageConstant and AddConstant", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{4}}}
		expr.AddDamageDice(1, 6, "Mace", tag.NewContainerFromString("bludgeoning"))
		expr.AddDamageConstant(2, "Enhancement", tag.NewContainer()) // Empty tags
		expr.AddConstant(1, "Feat Bonus")                            // Regular constant, empty tags

		grouped := expr.EvaluateGroup()

		assert.Len(t, grouped.Components, 1, "should have 1 damage type group")
		assert.Equal(t, 7, grouped.Components[0].Value, "total damage should be 4+2+1 = 7")
		assert.True(t, grouped.Components[0].Tags.HasAny(tag.NewContainerFromString("bludgeoning")))
	})
}

func TestExpression_ExpressionGroup(t *testing.T) {
	t.Run("can evaluate a constant", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 2, res.Value)
		assert.Len(t, res.Components, 1)
		assert.True(t, res.Components[0].Tags.HasAny(tag.NewContainerFromString("Slashing")), "expected Slashing tag to be present")
	})

	t.Run("groups primary under first component", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(3, "Str Mod", tag.NewContainerFromString("primary"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 5, res.Value)
		assert.Len(t, res.Components, 1)
		assert.True(t, res.Components[0].Tags.HasAny(tag.NewContainerFromString("Slashing")), "expected Slashing tag to be present")
	})

	t.Run("groups by tags", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(4, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(3, "Damage", tag.NewContainerFromString("Magical"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 9, res.Value)
		assert.Len(t, res.Components, 2)
		assert.True(t, res.Components[0].Tags.HasAny(tag.NewContainerFromString("Slashing")), "expected Slashing tag to be present")
		assert.True(t, res.Components[1].Tags.HasAny(tag.NewContainerFromString("Magical")), "expected Magical tag to be present")
	})

	t.Run("can evaluate a dice", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{5}}}
		expr.AddDamageDice(1, 6, "Damage", tag.NewContainerFromString("Slashing"))
		expr.AddDamageConstant(2, "Damage", tag.NewContainerFromString("Slashing"))
		res := expr.EvaluateGroup()
		assert.Equal(t, 7, res.Value)
		assert.Len(t, res.Components, 1)
		assert.True(t, res.Components[0].Tags.HasAny(tag.NewContainerFromString("Slashing")), "expected Slashing tag to be present")
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
