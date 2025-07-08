package expression

import (
	"testing"

	"anvil/internal/tag"

	"github.com/stretchr/testify/assert"
)

func TestExpression_New(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() Expression
		expected Component
	}{
		{
			name: "can create from constant",
			setup: func() Expression {
				return FromConstant(3, "Damage")
			},
			expected: Component{
				Type:   Constant,
				Source: "Damage",
				Value:  3,
			},
		},
		{
			name: "can create from dice",
			setup: func() Expression {
				return FromDice(2, 6, "Damage")
			},
			expected: Component{
				Type:   Dice,
				Source: "Damage",
				Sides:  6,
				Times:  2,
			},
		},
		{
			name: "can create from d20",
			setup: func() Expression {
				return FromD20("Damage")
			},
			expected: Component{
				Type:   D20,
				Source: "Damage",
				Sides:  20,
				Times:  1,
			},
		},
		{
			name: "can create from damage constant",
			setup: func() Expression {
				return FromDamageConstant(2, "Damage", tag.NewContainerFromString("slashing"))
			},
			expected: Component{
				Type:   DamageConstant,
				Source: "Damage",
				Tags:   tag.NewContainerFromString("slashing"),
				Value:  2,
			},
		},
		{
			name: "can create from damage dice",
			setup: func() Expression {
				return FromDamageDice(2, 6, "Damage", tag.NewContainerFromString("slashing"))
			},
			expected: Component{
				Type:   DamageDice,
				Source: "Damage",
				Tags:   tag.NewContainerFromString("slashing"),
				Times:  2,
				Sides:  6,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expression := tt.setup()
			assert.Equal(t, tt.expected.Type, expression.Components[0].Type)
			assert.Equal(t, tt.expected.Source, expression.Components[0].Source)
			if tt.expected.Value > 0 {
				assert.Equal(t, tt.expected.Value, expression.Components[0].Value)
			}
			if tt.expected.Times > 0 {
				assert.Equal(t, tt.expected.Times, expression.Components[0].Times)
			}
			if tt.expected.Sides > 0 {
				assert.Equal(t, tt.expected.Sides, expression.Components[0].Sides)
			}
			if !tt.expected.Tags.IsEmpty() {
				assert.True(t, expression.Components[0].Tags.HasAny(tt.expected.Tags), "expected tags to be present in component tags")
			}
		})
	}
}

func TestFromDamageResult(t *testing.T) {
	tests := []struct {
		name  string
		setup func() Expression
	}{
		{
			name: "creates clone of damage result",
			setup: func() Expression {
				expr := Expression{Value: 12}
				expr.AddDamageConstant(8, "Fire", tag.NewContainerFromString("fire"))
				expr.AddDamageConstant(4, "Cold", tag.NewContainerFromString("cold"))
				return expr
			},
		},
		{
			name: "creates clone of empty expression",
			setup: func() Expression {
				return Expression{Value: 0}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.setup()
			result := FromDamageResult(original)

			// Should be equal but not the same object
			assert.Equal(t, original.Value, result.Value)
			assert.Equal(t, len(original.Components), len(result.Components))

			// Verify it's a deep copy
			if len(result.Components) > 0 {
				result.Components[0].Value = 999
				assert.NotEqual(t, result.Components[0].Value, original.Components[0].Value)
			}
		})
	}
}
