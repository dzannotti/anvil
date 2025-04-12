package shapes

import "anvil/internal/grid"

func FloodFill(start grid.Position, radius int, isBlocked func(grid.Position) bool) []grid.Position {
	visited := make(map[grid.Position]bool)
	result := make([]grid.Position, 0)
	queue := []grid.Position{start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] || isBlocked(current) {
			continue
		}

		if current.Distance(start) > radius {
			continue
		}

		visited[current] = true
		result = append(result, current)

		neighbors := []grid.Position{
			{X: current.X + 1, Y: current.Y},
			{X: current.X - 1, Y: current.Y},
			{X: current.X, Y: current.Y + 1},
			{X: current.X, Y: current.Y - 1},
		}

		for _, next := range neighbors {
			if !visited[next] {
				queue = append(queue, next)
			}
		}
	}

	return result
}
