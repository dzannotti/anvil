package expression

import (
	"anvil/internal/tagcontainer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpression_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func() Expression
		expected Term
	}{
		{
			name: "can create from scalar",
			setup: func() Expression {
				return FromScalar(3, "Damage")
			},
			expected: Term{
				Type:   TypeScalar,
				Source: "Damage",
				Value:  3,
			},
		},
		{
			name: "can create from dice",
			setup: func() Expression {
				return FromDice(2, 6, "Damage")
			},
			expected: Term{
				Type:   TypeDice,
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
			expected: Term{
				Type:   TypeDice20,
				Source: "Damage",
				Sides:  20,
				Times:  1,
			},
		},
		{
			name: "can create from damage scalar",
			setup: func() Expression {
				return FromDamageScalar(2, "Damage", *tagcontainer.FromString("slashing"))
			},
			expected: Term{
				Type:   TypeDamageScalar,
				Source: "Damage",
				Tags:   tagcontainer.FromString("slashing"),
				Value:  2,
			},
		},
		{
			name: "can create from damage dice",
			setup: func() Expression {
				return FromDamageDice(2, 6, "Damage", *tagcontainer.FromString("slashing"))
			},
			expected: Term{
				Type:   TypeDamageDice,
				Source: "Damage",
				Tags:   tagcontainer.FromString("slashing"),
				Times:  2,
				Sides:  6,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expression := tt.setup()
			assert.Equal(t, tt.expected.Type, expression.Terms[0].Type)
			assert.Equal(t, tt.expected.Source, expression.Terms[0].Source)
			if tt.expected.Value > 0 {
				assert.Equal(t, tt.expected.Value, expression.Terms[0].Value)
			}
			if tt.expected.Times > 0 {
				assert.Equal(t, tt.expected.Times, expression.Terms[0].Times)
			}
			if tt.expected.Sides > 0 {
				assert.Equal(t, tt.expected.Sides, expression.Terms[0].Sides)
			}
			if tt.expected.Tags != nil {
				assert.Equal(t, tt.expected.Tags, expression.Terms[0].Tags)
			}
		})
	}
}
