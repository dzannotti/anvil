package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique_SliceElements(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "Integers with duplicates",
			input:    []int{1, 2, 2, 3, 3, 4, 5, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Strings with duplicates",
			input:    []string{"a", "b", "b", "c", "c", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "Single element",
			input:    []int{1},
			expected: []int{1},
		},
		{
			name:     "All duplicates",
			input:    []string{"test", "test", "test"},
			expected: []string{"test"},
		},
		{
			name:     "No duplicates",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				expected := tt.expected.([]int)
				result := SliceElements(input)
				if !assert.Equal(t, result, expected) {
					t.Errorf("SliceElements() = %v, want %v", result, expected)
				}
			case []string:
				expected := tt.expected.([]string)
				result := SliceElements(input)
				if !assert.Equal(t, result, expected) {
					t.Errorf("SliceElements() = %v, want %v", result, expected)
				}
			}
		})
	}
}
