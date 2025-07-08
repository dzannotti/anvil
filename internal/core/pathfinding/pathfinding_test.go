package pathfinding

import (
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/grid"
)

func TestBasicPathFinding(t *testing.T) {
	tests := []struct {
		name     string
		start    grid.Position
		end      grid.Position
		expected []grid.Position
	}{
		{
			name:  "straight horizontal path",
			start: grid.Position{X: 0, Y: 2},
			end:   grid.Position{X: 4, Y: 2},
			expected: []grid.Position{
				{X: 0, Y: 2},
				{X: 1, Y: 2},
				{X: 2, Y: 2},
				{X: 3, Y: 2},
				{X: 4, Y: 2},
			},
		},
		{
			name:  "straight vertical path",
			start: grid.Position{X: 2, Y: 0},
			end:   grid.Position{X: 2, Y: 4},
			expected: []grid.Position{
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 2, Y: 2},
				{X: 2, Y: 3},
				{X: 2, Y: 4},
			},
		},
		{
			name:  "diagonal path",
			start: grid.Position{X: 0, Y: 0},
			end:   grid.Position{X: 2, Y: 2},
			expected: []grid.Position{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			navCost := func(_ grid.Position) int {
				return 1
			}
			result, ok := FindPath(tt.start, tt.end, 5, 5, navCost)
			assert.True(t, ok, "path not found")

			assert.Equal(t, len(tt.expected), len(result.Path), "path length mismatch")

			for i := range result.Path {
				assert.Equal(t, tt.expected[i], result.Path[i], "path position %d mismatch", i)
			}
		})
	}
}

func TestObstacleAvoidance(t *testing.T) {
	t.Run("should navigate around obstacles", func(t *testing.T) {
		navCost := func(pos grid.Position) int {
			blocked := []grid.Position{
				{X: 2, Y: 1},
				{X: 2, Y: 2},
				{X: 2, Y: 3},
			}
			if slices.Contains(blocked, pos) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 1, Y: 2}
		end := grid.Position{X: 3, Y: 2}
		result, ok := FindPath(start, end, 5, 5, navCost)

		assert.True(t, ok, "path not found")

		assert.Greater(t, len(result.Path), 2, "path should be longer than direct route")

		assert.True(t, containsPosition(result.Path, start), "path should contain start position")
		assert.True(t, containsPosition(result.Path, end), "path should contain end position")

		for _, pos := range result.Path {
			cost := navCost(pos)
			assert.NotEqual(t, math.MaxInt, cost, "path contains obstacle at position %v", pos)
		}
	})

	t.Run("should return empty path when destination is unreachable", func(t *testing.T) {
		navCost := func(pos grid.Position) int {
			blocked := []grid.Position{}
			for y := range 5 {
				blocked = append(blocked, grid.Position{X: 2, Y: y})
			}
			if slices.Contains(blocked, pos) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 1, Y: 2}
		end := grid.Position{X: 3, Y: 2}
		result, ok := FindPath(start, end, 5, 5, navCost)

		assert.False(t, ok, "path found when not expected")

		assert.Nil(t, result, "path should be empty when destination is unreachable")
	})
}

func TestPathOptimality(t *testing.T) {
	t.Run("should prefer diagonal movement when it's shorter", func(t *testing.T) {
		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		navCost := func(_ grid.Position) int {
			return 1
		}
		result, ok := FindPath(start, end, 5, 5, navCost)

		assert.True(t, ok, "path not found")

		assert.Equal(t, 3, len(result.Path), "path should be diagonal with length 3")
	})

	t.Run("should find optimal path around obstacles", func(t *testing.T) {
		navCost := func(pos grid.Position) int {
			if pos == (grid.Position{X: 1, Y: 1}) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		result, ok := FindPath(start, end, 5, 5, navCost)

		assert.True(t, ok, "path not found")

		assert.Greater(t, len(result.Path), 3, "path should be longer than diagonal")
		assert.False(t, containsPosition(result.Path, grid.Position{X: 1, Y: 1}), "path should not contain obstacle position")
	})

	t.Run("cannot find diagonal path around walls", func(t *testing.T) {
		navCost := func(pos grid.Position) int {
			if pos == (grid.Position{X: 1, Y: 0}) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 0}
		result, ok := FindPath(start, end, 5, 5, navCost)

		assert.True(t, ok, "path not found")

		assert.Greater(t, len(result.Path), 3, "path should be longer than diagonal")
	})
}

func containsPosition(positions []grid.Position, target grid.Position) bool {
	for _, pos := range positions {
		if pos == target {
			return true
		}
	}
	return false
}
