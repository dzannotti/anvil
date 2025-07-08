package prettyprint

import (
	"testing"

	"anvil/internal/core"
	"anvil/internal/grid"

	"github.com/stretchr/testify/assert"
)

func TestPrintPosition(t *testing.T) {
	pos := grid.Position{X: 5, Y: 10}
	result := printPosition(pos)
	assert.Equal(t, "(5, 10)", result)
}

func TestPrintPositions(t *testing.T) {
	positions := []grid.Position{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 5, Y: 6},
	}
	result := printPositions(positions)
	assert.Equal(t, "(1, 2), (3, 4), (5, 6)", result)
}

func TestPrintPositions_Empty(t *testing.T) {
	result := printPositions([]grid.Position{})
	assert.Equal(t, "", result)
}

func TestPrintPositions_Single(t *testing.T) {
	positions := []grid.Position{{X: 7, Y: 8}}
	result := printPositions(positions)
	assert.Equal(t, "(7, 8)", result)
}

func TestPrintActorNames(t *testing.T) {
	actors := []*core.Actor{
		{Name: "Cedric"},
		{Name: "Zombie 1"},
		{Name: "Zombie 2"},
	}
	result := printActorNames(actors)
	assert.Equal(t, "Cedric, Zombie 1, Zombie 2", result)
}

func TestPrintActorNames_Empty(t *testing.T) {
	result := printActorNames([]*core.Actor{})
	assert.Equal(t, "", result)
}

func TestPrintActorNames_Single(t *testing.T) {
	actors := []*core.Actor{{Name: "Solo Actor"}}
	result := printActorNames(actors)
	assert.Equal(t, "Solo Actor", result)
}

func TestFormatRollResult(t *testing.T) {
	tests := []struct {
		name     string
		success  bool
		critical bool
		value    int
		against  int
		expected string
	}{
		{
			name:     "success",
			success:  true,
			critical: false,
			value:    15,
			against:  10,
			expected: "✅ Success 15 vs 10",
		},
		{
			name:     "failure",
			success:  false,
			critical: false,
			value:    8,
			against:  12,
			expected: "❌ Failure 8 vs 12",
		},
		{
			name:     "critical success",
			success:  true,
			critical: true,
			value:    20,
			against:  15,
			expected: "✅💥 Critical Success 20 vs 15",
		},
		{
			name:     "critical failure",
			success:  false,
			critical: true,
			value:    1,
			against:  10,
			expected: "❌💥 Critical Failure 1 vs 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRollResult(tt.success, tt.critical, tt.value, tt.against)
			assert.Equal(t, tt.expected, result)
		})
	}
}
