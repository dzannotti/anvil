package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"anvil/internal/grid"
)

func TestWorldIntegration(t *testing.T) {
	t.Run("World Components Work Together", func(t *testing.T) {
		t.Run("should create world with all components", func(t *testing.T) {
			world := NewWorld(10, 10)

			assert.Equal(t, 10, world.Width())
			assert.Equal(t, 10, world.Height())
			assert.NotNil(t, world.RequestManager())
			assert.NotNil(t, world.lineOfSightCalc)
		})

		t.Run("should handle spatial queries correctly", func(t *testing.T) {
			world := NewWorld(5, 5)
			actor := &Actor{Name: "TestActor"}
			pos := grid.Position{X: 2, Y: 2}

			world.AddOccupant(pos, actor)

			assert.Equal(t, actor, world.ActorAt(pos))
			assert.Nil(t, world.ActorAt(grid.Position{X: 0, Y: 0}))

			world.RemoveOccupant(pos, actor)
			assert.Nil(t, world.ActorAt(pos))
		})

		t.Run("should handle pathfinding correctly", func(t *testing.T) {
			world := NewWorld(5, 5)
			start := grid.Position{X: 0, Y: 0}
			end := grid.Position{X: 4, Y: 4}

			result, found := world.FindPath(start, end)

			assert.True(t, found)
			require.NotNil(t, result)
			assert.True(t, len(result.Steps) > 0)
		})

		t.Run("should handle line of sight correctly", func(t *testing.T) {
			world := NewWorld(5, 5)
			from := grid.Position{X: 0, Y: 0}
			to := grid.Position{X: 4, Y: 4}

			assert.True(t, world.HasLineOfSight(from, to))

			// Add a wall to block line of sight
			world.At(grid.Position{X: 2, Y: 2}).Tile = Wall
			assert.False(t, world.HasLineOfSight(from, to))
		})

		t.Run("should handle flood fill correctly", func(t *testing.T) {
			world := NewWorld(5, 5)
			start := grid.Position{X: 2, Y: 2}

			positions := world.FloodFill(start, 1)

			assert.True(t, len(positions) > 0)
			assert.Contains(t, positions, start)
		})

		t.Run("should handle actors in range correctly", func(t *testing.T) {
			world := NewWorld(5, 5)
			center := grid.Position{X: 2, Y: 2}

			actors := world.ActorsInRange(center, 1, func(_ *Actor) bool { return true })

			assert.Equal(t, 0, len(actors))
		})
	})

	t.Run("Request Manager Integration", func(t *testing.T) {
		t.Run("should access request manager through world", func(t *testing.T) {
			world := NewWorld(5, 5)

			assert.NotNil(t, world.RequestManager())
			assert.False(t, world.RequestManager().HasPendingRequest())
			assert.Nil(t, world.RequestManager().GetPendingRequest())
		})
	})

	t.Run("Error Handling", func(t *testing.T) {
		t.Run("should handle invalid positions", func(t *testing.T) {
			world := NewWorld(5, 5)

			assert.False(t, world.IsValidPosition(grid.Position{X: -1, Y: 0}))
			assert.False(t, world.IsValidPosition(grid.Position{X: 0, Y: -1}))
			assert.False(t, world.IsValidPosition(grid.Position{X: 5, Y: 0}))
			assert.False(t, world.IsValidPosition(grid.Position{X: 0, Y: 5}))
			assert.True(t, world.IsValidPosition(grid.Position{X: 0, Y: 0}))
			assert.True(t, world.IsValidPosition(grid.Position{X: 4, Y: 4}))
		})

		t.Run("should handle pathfinding to invalid positions", func(t *testing.T) {
			world := NewWorld(5, 5)
			start := grid.Position{X: 0, Y: 0}
			end := grid.Position{X: 10, Y: 10}

			result, found := world.FindPath(start, end)

			assert.False(t, found)
			require.NotNil(t, result)
		})

		t.Run("should handle pathfinding through walls", func(t *testing.T) {
			world := NewWorld(5, 5)
			// Create a wall barrier
			for x := 0; x < 5; x++ {
				world.At(grid.Position{X: x, Y: 2}).Tile = Wall
			}

			start := grid.Position{X: 0, Y: 0}
			end := grid.Position{X: 0, Y: 4}

			result, found := world.FindPath(start, end)

			assert.False(t, found)
			require.NotNil(t, result)
		})
	})
}
