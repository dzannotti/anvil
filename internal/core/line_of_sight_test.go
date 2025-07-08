package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/grid"
)

func TestLineOfSightCalculator(t *testing.T) {
	createTestWorld := func() *World {
		world := NewWorld(5, 5)
		// Create a simple 5x5 grid with some walls
		// Layout:
		// . . . . .
		// . W . . .
		// . W . . .
		// . . . . .
		// . . . . .
		world.At(grid.Position{X: 1, Y: 1}).Tile = Wall
		world.At(grid.Position{X: 1, Y: 2}).Tile = Wall
		return world
	}

	t.Run("Clear Line of Sight", func(t *testing.T) {
		t.Run("should have line of sight in open space", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 2, Y: 0}

			assert.True(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should have line of sight to adjacent cell", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 2, Y: 2}
			to := grid.Position{X: 3, Y: 2}

			assert.True(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should have line of sight to same position", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			pos := grid.Position{X: 2, Y: 2}

			assert.True(t, calc.HasLineOfSight(pos, pos))
		})
	})

	t.Run("Blocked Line of Sight", func(t *testing.T) {
		t.Run("should not have line of sight through wall", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 1}
			to := grid.Position{X: 2, Y: 1}

			assert.False(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should not have line of sight to wall", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 1}
			to := grid.Position{X: 1, Y: 1}

			assert.False(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should not have line of sight outside world bounds", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 10, Y: 10}

			assert.False(t, calc.HasLineOfSight(from, to))
		})
	})

	t.Run("Diagonal Line of Sight", func(t *testing.T) {

		t.Run("should be blocked by corner walls", func(t *testing.T) {
			world := NewWorld(5, 5)
			// Create a corner blocking scenario
			// . . . . .
			// . W W . .
			// . W W . .
			// . . . . .
			// . . . . .
			world.At(grid.Position{X: 1, Y: 1}).Tile = Wall
			world.At(grid.Position{X: 2, Y: 1}).Tile = Wall
			world.At(grid.Position{X: 1, Y: 2}).Tile = Wall
			world.At(grid.Position{X: 2, Y: 2}).Tile = Wall

			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 3, Y: 3}

			assert.False(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should allow diagonal movement around single wall", func(t *testing.T) {
			world := NewWorld(5, 5)
			// Create a single wall
			// . . . . .
			// . W . . .
			// . . . . .
			// . . . . .
			// . . . . .
			world.At(grid.Position{X: 1, Y: 1}).Tile = Wall

			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 2, Y: 1}

			assert.True(t, calc.HasLineOfSight(from, to))
		})
	})

	t.Run("Edge Cases", func(t *testing.T) {
		t.Run("should handle straight horizontal line", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 4, Y: 0}

			assert.True(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should handle straight vertical line", func(t *testing.T) {
			world := createTestWorld()
			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 0, Y: 4}

			assert.True(t, calc.HasLineOfSight(from, to))
		})

		t.Run("should handle line blocked by wall in middle", func(t *testing.T) {
			world := NewWorld(5, 5)
			// Create a wall in the middle of a horizontal line
			// . . W . .
			// . . . . .
			// . . . . .
			// . . . . .
			// . . . . .
			world.At(grid.Position{X: 2, Y: 0}).Tile = Wall

			calc := NewLineOfSightCalculator(world)

			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 4, Y: 0}

			assert.False(t, calc.HasLineOfSight(from, to))
		})
	})
}
