package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_Add(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() Expression
		expected struct {
			component  Component
			components []Component
		}
	}{
		{
			name: "can add a constant",
			setup: func() Expression {
				exp := Expression{}
				exp.AddConstant(2, "Damage")
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				component: Component{
					Type:   TypeConstant,
					Source: "Damage",
					Value:  2,
				},
			},
		},
		{
			name: "can add a dice",
			setup: func() Expression {
				exp := Expression{}
				exp.AddDice(2, 6, "Damage")
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				component: Component{
					Type:   TypeDice,
					Source: "Damage",
					Sides:  6,
					Times:  2,
				},
			},
		},
		{
			name: "can add a d20",
			setup: func() Expression {
				exp := Expression{}
				exp.AddD20("Damage")
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				component: Component{
					Type:   TypeDice20,
					Source: "Damage",
					Sides:  20,
					Times:  1,
				},
			},
		},
		{
			name: "can add a damage constant",
			setup: func() Expression {
				exp := Expression{}
				exp.AddDamageConstant(2, "Damage", tag.NewContainerFromString("slashing"))
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				component: Component{
					Type:   TypeDamageConstant,
					Source: "Damage",
					Tags:   tag.NewContainerFromString("slashing"),
					Value:  2,
				},
			},
		},
		{
			name: "can add a damage dice",
			setup: func() Expression {
				exp := Expression{}
				exp.AddDamageDice(2, 6, "Damage", tag.NewContainerFromString("slashing"))
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				component: Component{
					Type:   TypeDamageDice,
					Source: "Damage",
					Tags:   tag.NewContainerFromString("slashing"),
					Times:  2,
					Sides:  6,
				},
			},
		},
		{
			name: "can replace an expression with a component",
			setup: func() Expression {
				exp := Expression{}
				exp.AddConstant(2, "Damage")
				exp.ReplaceWith(3, "Bad Luck")
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				component: Component{
					Type:   TypeConstant,
					Source: "Bad Luck",
					Value:  3,
				},
			},
		},
		{
			name: "can add multiple components",
			setup: func() Expression {
				exp := Expression{}
				exp.AddConstant(2, "First")
				exp.AddDice(1, 6, "Second")
				return exp
			},
			expected: struct {
				component  Component
				components []Component
			}{
				components: []Component{
					{
						Type:   TypeConstant,
						Source: "First",
						Value:  2,
					},
					{
						Type:   TypeDice,
						Source: "Second",
						Times:  1,
						Sides:  6,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expression := tt.setup()

			if len(tt.expected.components) > 0 {
				assert.Len(t, expression.Components, len(tt.expected.components))
				for i, component := range tt.expected.components {
					assert.Equal(t, component.Type, expression.Components[i].Type)
					assert.Equal(t, component.Source, expression.Components[i].Source)
					if component.Value > 0 {
						assert.Equal(t, component.Value, expression.Components[i].Value)
					}
					if component.Times > 0 {
						assert.Equal(t, component.Times, expression.Components[i].Times)
					}
					if component.Sides > 0 {
						assert.Equal(t, component.Sides, expression.Components[i].Sides)
					}
				}
			} else {
				assert.Equal(t, tt.expected.component.Type, expression.Components[0].Type)
				assert.Equal(t, tt.expected.component.Source, expression.Components[0].Source)
				if tt.expected.component.Value > 0 {
					assert.Equal(t, tt.expected.component.Value, expression.Components[0].Value)
				}
				if tt.expected.component.Times > 0 {
					assert.Equal(t, tt.expected.component.Times, expression.Components[0].Times)
				}
				if tt.expected.component.Sides > 0 {
					assert.Equal(t, tt.expected.component.Sides, expression.Components[0].Sides)
				}
				if !tt.expected.component.Tags.IsEmpty() {
					assert.Equal(t, tt.expected.component.Tags, expression.Components[0].Tags)
				}
			}
		})
	}
}
