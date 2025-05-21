package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func() Expression
		expected struct {
			term  Term
			terms []Term
		}
	}{
		{
			name: "can add a scalar",
			setup: func() Expression {
				exp := Expression{}
				exp.AddScalar(2, "Damage")
				return exp
			},
			expected: struct {
				term  Term
				terms []Term
			}{
				term: Term{
					Type:   TypeScalar,
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
				term  Term
				terms []Term
			}{
				term: Term{
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
				term  Term
				terms []Term
			}{
				term: Term{
					Type:   TypeDice20,
					Source: "Damage",
					Sides:  20,
					Times:  1,
				},
			},
		},
		{
			name: "can add a damage scalar",
			setup: func() Expression {
				exp := Expression{}
				exp.AddDamageScalar(2, "Damage", tag.ContainerFromString("slashing"))
				return exp
			},
			expected: struct {
				term  Term
				terms []Term
			}{
				term: Term{
					Type:   TypeDamageScalar,
					Source: "Damage",
					Tags:   tag.ContainerFromString("slashing"),
					Value:  2,
				},
			},
		},
		{
			name: "can add a damage dice",
			setup: func() Expression {
				exp := Expression{}
				exp.AddDamageDice(2, 6, "Damage", tag.ContainerFromString("slashing"))
				return exp
			},
			expected: struct {
				term  Term
				terms []Term
			}{
				term: Term{
					Type:   TypeDamageDice,
					Source: "Damage",
					Tags:   tag.ContainerFromString("slashing"),
					Times:  2,
					Sides:  6,
				},
			},
		},
		{
			name: "can replace an expression with a term",
			setup: func() Expression {
				exp := Expression{}
				exp.AddScalar(2, "Damage")
				exp.ReplaceWith(3, "Bad Luck")
				return exp
			},
			expected: struct {
				term  Term
				terms []Term
			}{
				term: Term{
					Type:   TypeScalarReplace,
					Source: "Bad Luck",
					Value:  3,
				},
			},
		},
		{
			name: "can add multiple terms",
			setup: func() Expression {
				exp := Expression{}
				exp.AddScalar(2, "First")
				exp.AddDice(1, 6, "Second")
				return exp
			},
			expected: struct {
				term  Term
				terms []Term
			}{
				terms: []Term{
					{
						Type:   TypeScalar,
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

			if len(tt.expected.terms) > 0 {
				assert.Len(t, expression.Terms, len(tt.expected.terms))
				for i, term := range tt.expected.terms {
					assert.Equal(t, term.Type, expression.Terms[i].Type)
					assert.Equal(t, term.Source, expression.Terms[i].Source)
					if term.Value > 0 {
						assert.Equal(t, term.Value, expression.Terms[i].Value)
					}
					if term.Times > 0 {
						assert.Equal(t, term.Times, expression.Terms[i].Times)
					}
					if term.Sides > 0 {
						assert.Equal(t, term.Sides, expression.Terms[i].Sides)
					}
				}
			} else {
				assert.Equal(t, tt.expected.term.Type, expression.Terms[0].Type)
				assert.Equal(t, tt.expected.term.Source, expression.Terms[0].Source)
				if tt.expected.term.Value > 0 {
					assert.Equal(t, tt.expected.term.Value, expression.Terms[0].Value)
				}
				if tt.expected.term.Times > 0 {
					assert.Equal(t, tt.expected.term.Times, expression.Terms[0].Times)
				}
				if tt.expected.term.Sides > 0 {
					assert.Equal(t, tt.expected.term.Sides, expression.Terms[0].Sides)
				}
				if !tt.expected.term.Tags.IsEmpty() {
					assert.Equal(t, tt.expected.term.Tags, expression.Terms[0].Tags)
				}
			}
		})
	}
}
