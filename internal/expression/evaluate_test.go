package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpression_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() Expression
		expected struct {
			value     int
			component Component
			values    []int
			adv       []string
			dis       []string
		}
	}{
		{
			name: "can evaluate a constant",
			setup: func() Expression {
				exp := Expression{}
				exp.AddConstant(2, "Damage")
				return exp
			},
			expected: struct {
				value     int
				component Component
				values    []int
				adv       []string
				dis       []string
			}{
				value: 2,
				component: Component{
					Type:   TypeConstant,
					Source: "Damage",
				},
			},
		},
		{
			name: "can evaluate a dice",
			setup: func() Expression {
				exp := Expression{Rng: &mockRoller{mockReturns: []int{5}}}
				exp.AddDice(1, 6, "Damage")
				return exp
			},
			expected: struct {
				value     int
				component Component
				values    []int
				adv       []string
				dis       []string
			}{
				value: 5,
				component: Component{
					Type:   TypeDice,
					Source: "Damage",
					Sides:  6,
					Times:  1,
				},
				values: []int{5},
			},
		},
		{
			name: "can evaluate a simple d20",
			setup: func() Expression {
				exp := Expression{Rng: &mockRoller{mockReturns: []int{5}}}
				exp.AddD20("Attack")
				return exp
			},
			expected: struct {
				value     int
				component Component
				values    []int
				adv       []string
				dis       []string
			}{
				value: 5,
				component: Component{
					Type:   TypeDice20,
					Source: "Attack",
					Sides:  20,
					Times:  1,
				},
				values: []int{5},
			},
		},
		{
			name: "can evaluate d20 with advantage",
			setup: func() Expression {
				exp := Expression{Rng: &mockRoller{mockReturns: []int{5, 15}}}
				exp.AddD20("Attack")
				exp.WithAdvantage("Bless")
				return exp
			},
			expected: struct {
				value     int
				component Component
				values    []int
				adv       []string
				dis       []string
			}{
				value: 15,
				component: Component{
					Type:   TypeDice20,
					Source: "Attack",
				},
				values: []int{5, 15},
				adv:    []string{"Bless"},
			},
		},
		{
			name: "can evaluate d20 with disadvantage",
			setup: func() Expression {
				exp := Expression{Rng: &mockRoller{mockReturns: []int{5, 15}}}
				exp.AddD20("Attack")
				exp.WithDisadvantage("Cursed")
				return exp
			},
			expected: struct {
				value     int
				component Component
				values    []int
				adv       []string
				dis       []string
			}{
				value: 5,
				component: Component{
					Type:   TypeDice20,
					Source: "Attack",
				},
				values: []int{5, 15},
				dis:    []string{"Cursed"},
			},
		},
		{
			name: "can evaluate d20 with mixed advantage/disadvantage",
			setup: func() Expression {
				exp := Expression{Rng: &mockRoller{mockReturns: []int{5, 15}}}
				exp.AddD20("Attack")
				exp.WithAdvantage("Bless")
				exp.WithDisadvantage("Cursed")
				return exp
			},
			expected: struct {
				value     int
				component Component
				values    []int
				adv       []string
				dis       []string
			}{
				value: 5,
				component: Component{
					Type:   TypeDice20,
					Source: "Attack",
				},
				values: []int{5},
				adv:    []string{"Bless"},
				dis:    []string{"Cursed"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expression := tt.setup()
			res := expression.Evaluate()

			assert.Equal(t, tt.expected.value, res.Value)
			assert.Equal(t, tt.expected.component.Type, res.Components[0].Type)
			assert.Equal(t, tt.expected.component.Source, res.Components[0].Source)
			if tt.expected.component.Sides > 0 {
				assert.Equal(t, tt.expected.component.Sides, res.Components[0].Sides)
			}
			if tt.expected.component.Times > 0 {
				assert.Equal(t, tt.expected.component.Times, res.Components[0].Times)
			}
			if len(tt.expected.values) > 0 {
				assert.Equal(t, tt.expected.values, res.Components[0].Values)
			}
			if len(tt.expected.adv) > 0 {
				assert.Equal(t, tt.expected.adv, res.Components[0].HasAdvantage)
			}
			if len(tt.expected.dis) > 0 {
				assert.Equal(t, tt.expected.dis, res.Components[0].HasDisadvantage)
			}
		})
	}
}
