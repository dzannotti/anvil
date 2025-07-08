package pathfinding

import (
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			require := require.New(t)
			assert := assert.New(t)

			navCost := func(_ grid.Position) int { return 1 }
			result := FindPath(tt.start, tt.end, 5, 5, navCost)

			require.True(result.Found, "path should be found")
			require.Len(result.Steps, len(tt.expected), "path length should match expected")

			for i := range result.Steps {
				assert.Equal(tt.expected[i], result.Steps[i].Position, "path position at step %d", i)
			}
		})
	}
}

func TestObstacleAvoidance(t *testing.T) {
	t.Run("should navigate around obstacles", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			blocked := []grid.Position{
				{X: 2, Y: 1}, {X: 2, Y: 2}, {X: 2, Y: 3},
			}
			if slices.Contains(blocked, pos) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 1, Y: 2}
		end := grid.Position{X: 3, Y: 2}
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "path should be found")
		assert.Greater(len(result.Steps), 2, "path should be longer than direct route")
		assert.True(containsStep(result.Steps, start), "path should contain start position")
		assert.True(containsStep(result.Steps, end), "path should contain end position")

		for _, step := range result.Steps {
			cost := navCost(step.Position)
			assert.NotEqual(math.MaxInt, cost, "path should not contain obstacle at position %v", step.Position)
		}
	})

	t.Run("should return empty path when destination is unreachable", func(t *testing.T) {
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			// Block entire column
			if pos.X == 2 {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 1, Y: 2}
		end := grid.Position{X: 3, Y: 2}
		result := FindPath(start, end, 5, 5, navCost)

		assert.False(result.Found, "path should not be found")
		assert.Empty(result.Steps, "steps should be empty when destination is unreachable")
	})
}

func TestPathOptimality(t *testing.T) {
	t.Run("should prefer diagonal movement when it's shorter", func(t *testing.T) {
		require := require.New(t)

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		navCost := func(_ grid.Position) int { return 1 }
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "path should be found")
		require.Len(result.Steps, 3, "path should be diagonal with length 3")
	})

	t.Run("should find optimal path around obstacles", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			if pos == (grid.Position{X: 1, Y: 1}) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "path should be found")
		assert.Greater(len(result.Steps), 3, "path should be longer than diagonal")
		assert.False(containsStep(result.Steps, grid.Position{X: 1, Y: 1}), "path should not contain obstacle position")
	})

	t.Run("cannot find diagonal path around walls", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			if pos == (grid.Position{X: 1, Y: 0}) {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 0}
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "path should be found")
		assert.Greater(len(result.Steps), 3, "path should be longer than direct")
	})
}

func TestRichMetadata(t *testing.T) {
	t.Run("should provide rich metadata for each step", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 0}
		navCost := func(_ grid.Position) int { return 1 }
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "path should be found")
		require.Len(result.Steps, 3, "expected 3 steps")

		// Check first step
		assert.Equal(start, result.Steps[0].Position, "first step should be start position")
		assert.Equal(0.0, result.Steps[0].GCost, "first step GCost should be 0")
		assert.Equal(0.0, result.Steps[0].StepCost, "first step StepCost should be 0")
		assert.Equal(0, result.Steps[0].Distance, "first step Distance should be 0")

		// Check second step
		expectedSecond := grid.Position{X: 1, Y: 0}
		assert.Equal(expectedSecond, result.Steps[1].Position, "second step position")
		assert.Equal(1.0, result.Steps[1].GCost, "second step GCost should be 1.0")
		assert.Equal(1.0, result.Steps[1].StepCost, "second step StepCost should be 1.0")
		assert.Equal(1, result.Steps[1].Distance, "second step Distance should be 1")

		// Check total cost
		assert.Equal(2.0, result.TotalCost, "total cost should be 2.0")
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("single position path", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		start := grid.Position{X: 2, Y: 2}
		end := grid.Position{X: 2, Y: 2}
		navCost := func(_ grid.Position) int { return 1 }

		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "should find path to same position")
		require.Len(result.Steps, 1, "single position should have 1 step")
		assert.Equal(start, result.Steps[0].Position, "single step should be the start position")
		assert.Equal(0.0, result.TotalCost, "cost to same position should be 0")
	})

	t.Run("path along grid boundary", func(t *testing.T) {
		require := require.New(t)

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 4, Y: 0}
		navCost := func(_ grid.Position) int { return 1 }

		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "should find path along boundary")
		require.Len(result.Steps, 5, "boundary path should have 5 steps")
	})

	t.Run("blocked destination", func(t *testing.T) {
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			if pos == (grid.Position{X: 2, Y: 2}) {
				return math.MaxInt // Destination is blocked
			}
			return 1
		}

		result := FindPath(grid.Position{X: 0, Y: 0}, grid.Position{X: 2, Y: 2}, 5, 5, navCost)

		assert.False(result.Found, "should not find path when destination is blocked")
	})
}

func TestVariableCosts(t *testing.T) {
	t.Run("should prefer lower cost paths", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			// Make middle column expensive
			if pos.X == 1 {
				return 5
			}
			return 1
		}

		start := grid.Position{X: 0, Y: 1}
		end := grid.Position{X: 2, Y: 1}
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "should find path")

		// Should avoid the expensive middle column
		for _, step := range result.Steps {
			if step.Position.X == 1 {
				// If it goes through expensive area, total cost should reflect that
				assert.GreaterOrEqual(result.TotalCost, 3.0, "should account for high movement costs")
			}
		}
	})

	t.Run("should handle mixed terrain costs", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(pos grid.Position) int {
			// Create different terrain types
			switch pos.X {
			case 0, 4:
				return 1 // Normal
			case 1, 3:
				return 2 // Difficult
			case 2:
				return 3 // Very difficult
			default:
				return 1
			}
		}

		start := grid.Position{X: 0, Y: 2}
		end := grid.Position{X: 4, Y: 2}
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "should find path through mixed terrain")

		// Verify costs accumulate correctly
		expectedMinCost := float64(1 + 2 + 3 + 2) // Base costs
		assert.GreaterOrEqual(result.TotalCost, expectedMinCost, "total cost should account for terrain")
	})
}

func TestComplexScenarios(t *testing.T) {
	t.Run("maze with single solution", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		// Create a simple maze:
		// S.###
		// #...#
		// ###.E
		navCost := func(pos grid.Position) int {
			walls := map[grid.Position]bool{
				{X: 2, Y: 0}: true, {X: 3, Y: 0}: true, {X: 4, Y: 0}: true,
				{X: 0, Y: 1}: true, {X: 4, Y: 1}: true,
				{X: 0, Y: 2}: true, {X: 1, Y: 2}: true, {X: 2, Y: 2}: true,
			}
			if walls[pos] {
				return math.MaxInt
			}
			return 1
		}

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 4, Y: 2}
		result := FindPath(start, end, 5, 3, navCost)

		require.True(result.Found, "should find path through maze")

		// Verify no path goes through walls
		for _, step := range result.Steps {
			assert.NotEqual(math.MaxInt, navCost(step.Position), "path should not go through wall at %v", step.Position)
		}

		// Should be forced to take the long way around
		assert.GreaterOrEqual(len(result.Steps), 7, "path should be long enough for this maze")
	})

	t.Run("multiple equivalent paths", func(t *testing.T) {
		require := require.New(t)

		// Simple open area where multiple paths have same cost
		navCost := func(_ grid.Position) int { return 1 }

		start := grid.Position{X: 0, Y: 0}
		end := grid.Position{X: 2, Y: 2}
		result := FindPath(start, end, 5, 5, navCost)

		require.True(result.Found, "should find path")

		// Should find optimal diagonal path
		require.Len(result.Steps, 3, "should find optimal 3-step diagonal path")
	})
}

func TestCostAccuracy(t *testing.T) {
	t.Run("diagonal vs straight cost comparison", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(_ grid.Position) int { return 1 }

		// Test straight path
		straightResult := FindPath(grid.Position{X: 0, Y: 0}, grid.Position{X: 3, Y: 0}, 5, 5, navCost)
		// Test diagonal path
		diagonalResult := FindPath(grid.Position{X: 0, Y: 0}, grid.Position{X: 2, Y: 2}, 5, 5, navCost)

		require.True(straightResult.Found, "straight path should be found")
		require.True(diagonalResult.Found, "diagonal path should be found")

		// Diagonal should cost more per step due to 1.4x multiplier
		straightCostPerStep := straightResult.TotalCost / float64(len(straightResult.Steps)-1)
		diagonalCostPerStep := diagonalResult.TotalCost / float64(len(diagonalResult.Steps)-1)

		assert.Greater(diagonalCostPerStep, straightCostPerStep, "diagonal cost per step should be higher than straight")
	})

	t.Run("step cost calculation accuracy", func(t *testing.T) {
		require := require.New(t)
		assert := assert.New(t)

		navCost := func(_ grid.Position) int { return 1 }

		result := FindPath(grid.Position{X: 0, Y: 0}, grid.Position{X: 1, Y: 1}, 5, 5, navCost)

		require.True(result.Found, "path should be found")
		require.Len(result.Steps, 2, "diagonal path should have 2 steps")

		// First step should have no cost
		assert.Equal(0.0, result.Steps[0].StepCost, "first step should have 0 cost")

		// Second step should be diagonal cost (1.4)
		assert.Equal(1.4, result.Steps[1].StepCost, "diagonal step should cost 1.4")
	})
}

// Helper functions
func containsStep(steps []PathStep, target grid.Position) bool {
	for _, step := range steps {
		if step.Position == target {
			return true
		}
	}
	return false
}
