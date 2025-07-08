package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestExpression_Clone(t *testing.T) {
	tests := []struct {
		name  string
		setup func() Expression
	}{
		{
			name: "clones simple expression",
			setup: func() Expression {
				expr := Expression{Value: 42}
				expr.AddConstant(10, "Base")
				return expr
			},
		},
		{
			name: "clones complex expression with multiple components",
			setup: func() Expression {
				expr := Expression{Value: 15, Rng: &mockRoller{}}
				expr.AddD20("Attack")
				expr.AddConstant(5, "Bonus")
				expr.AddDamageDice(2, 6, "Damage", tag.NewContainerFromString("fire"))
				return expr
			},
		},
		{
			name: "clones empty expression",
			setup: func() Expression {
				return Expression{Value: 0}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.setup()
			clone := original.Clone()

			// Verify values are equal
			assert.Equal(t, original.Value, clone.Value)
			assert.Equal(t, original.Rng, clone.Rng)
			assert.Equal(t, len(original.Components), len(clone.Components))

			// Verify components are deeply cloned
			for i := range original.Components {
				assert.Equal(t, original.Components[i], clone.Components[i])

				// Verify it's a deep copy - modifying clone shouldn't affect original
				if len(clone.Components) > 0 {
					clone.Components[i].Value = 999
					assert.NotEqual(t, clone.Components[i].Value, original.Components[i].Value)
				}
			}

			// Verify modifying clone doesn't affect original
			clone.Value = 999
			assert.NotEqual(t, clone.Value, original.Value)
		})
	}
}
