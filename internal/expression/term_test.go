package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerm_ShouldModifyRoll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		term     Term
		expected bool
	}{
		{
			name:     "should return false when no modifiers are present",
			term:     Term{},
			expected: false,
		},
		{
			name: "should return true when advantage is present",
			term: Term{
				HasAdvantage: []string{"test"},
			},
			expected: true,
		},
		{
			name: "should return true when disadvantage is present",
			term: Term{
				HasDisadvantage: []string{"test"},
			},
			expected: true,
		},
		{
			name: "should return false when both advantage and disadvantage are present",
			term: Term{
				HasAdvantage:    []string{"test1", "test3"},
				HasDisadvantage: []string{"test2"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.term.shouldModifyRoll()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTerm_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		originalTerm Term
		modifyFunc   func(*Term)
		verifyFunc   func(*testing.T, Term, Term)
	}{
		{
			name: "should deep copy all fields and slices",
			originalTerm: func() Term {
				term := makeTerm(TypeDice, "test")
				term.Value = 10
				term.Values = []int{1, 2, 3}
				term.Times = 2
				term.Sides = 6
				term.HasAdvantage = []string{"adv1", "adv2"}
				term.HasDisadvantage = []string{"dis1"}
				term.Terms = []Term{makeTerm(TypeScalar, "sub")}
				return term
			}(),
			modifyFunc: func(t *Term) {
				t.Values[0] = 99
				t.HasAdvantage[0] = "modified"
				t.HasDisadvantage[0] = "modified"
			},
			verifyFunc: func(t *testing.T, original Term, clone Term) {
				// Verify all fields are equal
				assert.Equal(t, original.Type, clone.Type)
				assert.Equal(t, original.Value, clone.Value)
				assert.Equal(t, original.Source, clone.Source)
				assert.Equal(t, original.Times, clone.Times)
				assert.Equal(t, original.Sides, clone.Sides)

				// Verify slices are deep copied
				assert.Equal(t, original.Values, clone.Values)
				assert.NotSame(t, &original.Values, &clone.Values)

				assert.Equal(t, original.HasAdvantage, clone.HasAdvantage)
				assert.NotSame(t, &original.HasAdvantage, &clone.HasAdvantage)

				assert.Equal(t, original.HasDisadvantage, clone.HasDisadvantage)
				assert.NotSame(t, &original.HasDisadvantage, &clone.HasDisadvantage)

				// Verify nested Terms are cloned
				assert.Equal(t, original.Terms, clone.Terms)
				assert.NotSame(t, &original.Terms, &clone.Terms)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.originalTerm.Clone()
			tt.verifyFunc(t, tt.originalTerm, clone)

			// Modify original to verify deep copy
			tt.modifyFunc(&tt.originalTerm)

			// Verify modifications don't affect clone
			assert.NotEqual(t, tt.originalTerm.Values[0], clone.Values[0])
			assert.NotEqual(t, tt.originalTerm.HasAdvantage[0], clone.HasAdvantage[0])
			assert.NotEqual(t, tt.originalTerm.HasDisadvantage[0], clone.HasDisadvantage[0])
		})
	}
}
