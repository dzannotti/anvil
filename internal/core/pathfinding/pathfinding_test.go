package pathfinding

import (
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
			pathFinding := New(5, 5)
			result, ok := pathFinding.FindPath(tt.start, tt.end)
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
		pathFinding := New(5, 5)
		obstacles := make([][]bool, 5)
		for i := range obstacles {
			obstacles[i] = make([]bool, 5)
		}
		node, ok := pathFinding.At(grid.Position{X: 2, Y: 1})
		if ok {
			node.Walkable = false
		}
		node, ok = pathFinding.At(grid.Position{X: 2, Y: 2})
		if ok {
			node.Walkable = false
		}
		node, ok = pathFinding.At(grid.Position{X: 2, Y: 3})
		if ok {
			node.Walkable = false
		}

		start := grid.Position{X: 1, Y: 2}
		end := grid.Position{X: 3, Y: 2}
		result, ok := pathFinding.FindPath(start, end)

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
			node, _ := pathFinding.At(pos)
			if !node.Walkable {
				t.Errorf("path contains obstacle at position %v", pos)
			}
		}
	})

	t.Run("should return empty path when destination is unreachable", func(t *testing.T) {
		pathFinding := New(5, 5)
		obstacles := make([][]bool, 5)
		for i := range obstacles {
			obstacles[i] = make([]bool, 5)
		}
		for y := 0; y < 5; y++ {
			node, ok := pathFinding.At(grid.Position{X: 2, Y: y})
			if ok {
				node.Walkable = false
			}
		}

		start := grid.Position{X: 1, Y: 2}
		end := grid.Position{X: 3, Y: 2}
		result, ok := pathFinding.FindPath(start, end)

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
		pathFinding := New(5, 5)
		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		result, ok := pathFinding.FindPath(start, end)

		if !ok {
			t.Error("path not found")
		}

		if len(result.Path) != 3 {
			t.Error("path should be diagonal with length 3")
		}
	})

	t.Run("should find optimal path around obstacles", func(t *testing.T) {
		pathFinding := New(5, 5)
		obstacles := make([][]bool, 5)
		for i := range obstacles {
			obstacles[i] = make([]bool, 5)
		}
		node, ok := pathFinding.At(grid.Position{X: 1, Y: 1})
		if ok {
			node.Walkable = false
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		result, _ := pathFinding.FindPath(start, end)

		if len(result.Path) <= 3 {
			t.Error("path should be longer than diagonal")
		}
		if containsPosition(result.Path, grid.Position{X: 1, Y: 1}) {
			t.Error("path should not contain obstacle position")
		}
	})

	t.Run("cannot find diagonal path around walls", func(t *testing.T) {
		pathFinding := New(5, 5)
		obstacles := make([][]bool, 5)
		for i := range obstacles {
			obstacles[i] = make([]bool, 5)
		}
		node, ok := pathFinding.At(grid.Position{X: 1, Y: 0})
		if ok {
			node.Walkable = false
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 0}
		result, _ := pathFinding.FindPath(start, end)

		if len(result.Path) <= 3 {
			t.Error("path should be longer than diagonal")
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("should calculate correct distance", func(t *testing.T) {
		pathFinding := New(5, 5)
		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 3, Y: 3}
		distance := pathFinding.distance(start, end)

		if distance != 3 {
			t.Errorf("distance = %v, want %v", distance, 3)
		}
	})

	t.Run("should generate valid neighbor nodes", func(t *testing.T) {
		pathFinding := New(5, 5)
		center := grid.Position{X: 2, Y: 2}
		node, _ := pathFinding.grid.At(center)
		neighbors := pathFinding.neighbours(node)

		if len(neighbors) != 8 {
			t.Errorf("number of neighbors = %v, want 8", len(neighbors))
		}

		for _, node := range neighbors {
			pos := node.Position
			if pos.X < 0 || pos.X >= 5 || pos.Y < 0 || pos.Y >= 5 {
				t.Errorf("neighbor position %v is out of bounds", pos)
			}
		}
	})

	t.Run("should generate fewer neighbors at edges", func(t *testing.T) {
		pathFinding := New(5, 5)
		corner := grid.Position{X: 0, Y: 0}
		node, _ := pathFinding.grid.At(corner)
		neighbors := pathFinding.neighbours(node)

		if len(neighbors) != 3 {
			t.Errorf("number of neighbors = %v, want 3", len(neighbors))
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
