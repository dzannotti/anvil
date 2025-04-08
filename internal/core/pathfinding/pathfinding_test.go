package pathfinding

import (
	"math"
	"slices"
	"testing"

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
			navCost := func(pos grid.Position) int {
				return 1
			}
			result, ok := FindPath(tt.start, tt.end, 5, 5, navCost)
			if !ok {
				t.Errorf("path not found")
			}

			if len(result.Path) != len(tt.expected) {
				t.Errorf("path length = %v, want %v", len(result.Path), len(tt.expected))
			}

			for i := range result.Path {
				if result.Path[i] != tt.expected[i] {
					t.Errorf("path[%d] = %v, want %v", i, result.Path[i], tt.expected[i])
				}
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

		if !ok {
			t.Error("path not found")
		}

		if len(result.Path) <= 2 {
			t.Error("path should be longer than direct route")
		}

		if !containsPosition(result.Path, start) {
			t.Error("path should contain start position")
		}
		if !containsPosition(result.Path, end) {
			t.Error("path should contain end position")
		}

		for _, pos := range result.Path {
			cost := navCost(pos)
			if cost == math.MaxInt {
				t.Errorf("path contains obstacle at position %v", pos)
			}
		}
	})

	t.Run("should return empty path when destination is unreachable", func(t *testing.T) {
		navCost := func(pos grid.Position) int {
			blocked := []grid.Position{}
			for y := 0; y < 5; y++ {
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

		if ok {
			t.Error("path found when not expected")
		}

		if result != nil {
			t.Error("path should be empty when destination is unreachable")
		}
	})
}

func TestPathOptimality(t *testing.T) {
	t.Run("should prefer diagonal movement when it's shorter", func(t *testing.T) {
		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		navCost := func(pos grid.Position) int {
			return 1
		}
		result, ok := FindPath(start, end, 5, 5, navCost)

		if !ok {
			t.Error("path not found")
		}

		if len(result.Path) != 3 {
			t.Error("path should be diagonal with length 3")
		}
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

		if !ok {
			t.Error("path not found")
		}

		if len(result.Path) <= 3 {
			t.Error("path should be longer than diagonal")
		}
		if containsPosition(result.Path, grid.Position{X: 1, Y: 1}) {
			t.Error("path should not contain obstacle position")
		}
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

		if !ok {
			t.Error("path not found")
		}

		if len(result.Path) <= 3 {
			t.Error("path should be longer than diagonal")
		}
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
