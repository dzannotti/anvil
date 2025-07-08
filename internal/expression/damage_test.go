package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_HalveDamage(t *testing.T) {
	t.Run("halves matching damage component", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{6, 6}}}
		expr.AddDamageDice(1, 6, "Fire Damage", tag.NewContainerFromString("fire"))
		expr.Evaluate()

		expr.HalveDamage(tag.FromString("fire"), "Resistance")

		assert.Equal(t, 3, expr.Components[0].Value) // floor(6/2)
		assert.Equal(t, "Halved (Resistance) Fire Damage", expr.Components[0].Source)
		assert.Equal(t, Constant, expr.Components[0].Type)
	})

	t.Run("ignores non-matching damage components", func(t *testing.T) {
		expr := Expression{Rng: &mockRoller{mockReturns: []int{8, 8}}}
		expr.AddDamageDice(1, 8, "Cold Damage", tag.NewContainerFromString("cold"))
		expr.Evaluate()

		expr.HalveDamage(tag.FromString("fire"), "Resistance")

		// Component should remain unchanged since tag doesn't match
		assert.Equal(t, 8, expr.Components[0].Value)
		assert.Equal(t, "Cold Damage", expr.Components[0].Source)
		assert.Equal(t, DamageDice, expr.Components[0].Type) // Should remain as dice
	})

	t.Run("handles odd numbers by rounding down", func(t *testing.T) {
		expr := Expression{}
		expr.AddDamageConstant(7, "Poison Damage", tag.NewContainerFromString("poison"))

		expr.HalveDamage(tag.FromString("poison"), "Resistance")

		assert.Equal(t, 3, expr.Components[0].Value) // floor(7/2)
		assert.Equal(t, "Halved (Resistance) Poison Damage", expr.Components[0].Source)
		assert.Equal(t, Constant, expr.Components[0].Type)
	})
}

func TestExpression_HasDamageType(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() Expression
		checkTag tag.Tag
		expected bool
	}{
		{
			name: "returns true when damage type is present",
			setup: func() Expression {
				expr := Expression{}
				expr.AddDamageConstant(5, "Fire Damage", tag.NewContainerFromString("fire"))
				return expr
			},
			checkTag: tag.FromString("fire"),
			expected: true,
		},
		{
			name: "returns false when damage type is not present",
			setup: func() Expression {
				expr := Expression{}
				expr.AddDamageConstant(5, "Fire Damage", tag.NewContainerFromString("fire"))
				return expr
			},
			checkTag: tag.FromString("cold"),
			expected: false,
		},
		{
			name: "returns true when one of multiple components matches",
			setup: func() Expression {
				expr := Expression{}
				expr.AddDamageConstant(5, "Fire Damage", tag.NewContainerFromString("fire"))
				expr.AddDamageConstant(3, "Cold Damage", tag.NewContainerFromString("cold"))
				return expr
			},
			checkTag: tag.FromString("cold"),
			expected: true,
		},
		{
			name: "returns false for empty expression",
			setup: func() Expression {
				return Expression{}
			},
			checkTag: tag.FromString("fire"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.setup()
			result := expr.HasDamageType(tt.checkTag)
			assert.Equal(t, tt.expected, result)
		})
	}
}
