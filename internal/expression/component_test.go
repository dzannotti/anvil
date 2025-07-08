package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponent_ShouldModifyRoll(t *testing.T) {
	tests := []struct {
		name      string
		component Component
		expected  bool
	}{
		{
			name:      "should return false when no modifiers are present",
			component: Component{},
			expected:  false,
		},
		{
			name: "should return true when advantage is present",
			component: Component{
				HasAdvantage: []string{"test"},
			},
			expected: true,
		},
		{
			name: "should return true when disadvantage is present",
			component: Component{
				HasDisadvantage: []string{"test"},
			},
			expected: true,
		},
		{
			name: "should return false when both advantage and disadvantage are present",
			component: Component{
				HasAdvantage:    []string{"test1", "test3"},
				HasDisadvantage: []string{"test2"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.component.hasRollModifier()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComponent_Clone(t *testing.T) {
	tests := []struct {
		name              string
		originalComponent Component
		modifyFunc        func(*Component)
		verifyFunc        func(*testing.T, Component, Component)
	}{
		{
			name: "should deep copy all fields and slices",
			originalComponent: Component{
				Type:            Dice,
				Source:          "test",
				Value:           10,
				Values:          []int{1, 2, 3},
				Times:           2,
				Sides:           6,
				HasAdvantage:    []string{"adv1", "adv2"},
				HasDisadvantage: []string{"dis1"},
				Components: []Component{{
					Type:   Constant,
					Source: "sub",
				}},
			},
			modifyFunc: func(c *Component) {
				c.Values[0] = 99
				c.HasAdvantage[0] = "modified"
				c.HasDisadvantage[0] = "modified"
			},
			verifyFunc: func(t *testing.T, original Component, clone Component) {
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

				// Verify nested Components are cloned
				assert.Equal(t, original.Components, clone.Components)
				assert.NotSame(t, &original.Components, &clone.Components)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.originalComponent.Clone()
			tt.verifyFunc(t, tt.originalComponent, clone)

			// Modify original to verify deep copy
			tt.modifyFunc(&tt.originalComponent)

			// Verify modifications don't affect clone
			assert.NotEqual(t, tt.originalComponent.Values[0], clone.Values[0])
			assert.NotEqual(t, tt.originalComponent.HasAdvantage[0], clone.HasAdvantage[0])
			assert.NotEqual(t, tt.originalComponent.HasDisadvantage[0], clone.HasDisadvantage[0])
		})
	}
}
