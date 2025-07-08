package shapes

import "anvil/internal/grid"

//nolint:gocognit,cyclop // reason: cyclop here is allowed
func FloodFill(start grid.Position, radius int, isBlocked func(grid.Position) bool) []grid.Position {
	visited := make(map[grid.Position]bool)
	result := make([]grid.Position, 0)
	queue := []grid.Position{start}

	directions := []grid.Position{
		{X: 1, Y: 0},  // Right
		{X: -1, Y: 0}, // Left
		{X: 0, Y: 1},  // Down
		{X: 0, Y: -1}, // Up
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		if isBlocked(current) {
			continue
		}

		if current.Distance(start) > radius {
			continue
		}

		visited[current] = true
		result = append(result, current)

		for _, dir := range directions {
			next := grid.Position{X: current.X + dir.X, Y: current.Y + dir.Y}

			if !visited[next] && !isBlocked(next) {
				queue = append(queue, next)
			}
		}
	}

	return result
}
