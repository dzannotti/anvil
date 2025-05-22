package pathfinding

import (
	"math"
	"slices"

	"anvil/internal/grid"
	"anvil/internal/mathi"
)

type Result struct {
	Path  []grid.Position
	Cost  int
	Speed int
}

type node struct {
	pos    grid.Position
	fScore int
}

var offsets = []grid.Position{
	{X: 0, Y: -1},  // up
	{X: 0, Y: 1},   // down
	{X: -1, Y: 0},  // left
	{X: 1, Y: 0},   // right
	{X: -1, Y: -1}, // up-left
	{X: 1, Y: -1},  // up-right
	{X: -1, Y: 1},  // down-left
	{X: 1, Y: 1},   // down-right
}

//nolint:gocognit,cyclop // reason: cyclop here is allowed
func FindPath(
	start grid.Position,
	end grid.Position,
	width int,
	height int,
	movementCost func(grid.Position) int,
) (*Result, bool) {
	size := width * height
	gCost := make([]int, size)
	cameFrom := make([]*grid.Position, size)

	for i := range gCost {
		gCost[i] = math.MaxInt
	}

	idx := func(p grid.Position) int {
		return p.Y*width + p.X
	}

	inBounds := func(p grid.Position) bool {
		return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
	}

	open := newMinHeap()

	startNode := &node{
		pos:    start,
		fScore: heuristic(start, end),
	}
	open.Push(startNode)
	gCost[idx(start)] = 0

	for !open.Empty() {
		current := open.Pop()
		currentIdx := idx(current.pos)

		if current.pos == end {
			return reconstructPath(cameFrom, end, width, gCost[idx(end)]), true
		}

		for i, offset := range offsets {
			neighborPos := current.pos.Add(offset)
			if !inBounds(neighborPos) {
				continue
			}

			cost := movementCost(neighborPos)
			if cost == math.MaxInt {
				continue
			}

			// Check diagonal wall cutting only for diagonal moves
			isDiagonal := i >= 4
			if isDiagonal {
				dx := offset.X
				dy := offset.Y
				adj1 := grid.Position{X: current.pos.X + dx, Y: current.pos.Y}
				adj2 := grid.Position{X: current.pos.X, Y: current.pos.Y + dy}
				if movementCost(adj1) == math.MaxInt || movementCost(adj2) == math.MaxInt {
					continue
				}
			}

			moveCost := cost * 10
			if isDiagonal {
				moveCost += 4
			}

			neighborIdx := idx(neighborPos)
			tentativeG := gCost[currentIdx] + moveCost

			if tentativeG < gCost[neighborIdx] {
				gCost[neighborIdx] = tentativeG
				if isDiagonal {
					gCost[neighborIdx]++
				}

				cameFrom[neighborIdx] = &current.pos

				open.Push(&node{
					pos:    neighborPos,
					fScore: tentativeG + heuristic(neighborPos, end),
				})
			}
		}
	}

	return nil, false
}

func reconstructPath(cameFrom []*grid.Position, end grid.Position, width int, cost int) *Result {
	path := []grid.Position{end}
	curr := end

	idx := func(p grid.Position) int {
		return p.Y*width + p.X
	}

	for {
		prevPtr := cameFrom[idx(curr)]
		if prevPtr == nil {
			break
		}
		curr = *prevPtr
		path = append(path, curr)
	}

	slices.Reverse(path)
	return &Result{
		Path:  path,
		Cost:  cost,
		Speed: cost / 10,
	}
}
func heuristic(a, b grid.Position) int {
	dx := mathi.Abs(a.X - b.X)
	dy := mathi.Abs(a.Y - b.Y)
	// Chebyshev distance - maximum of dx and dy
	return 10 * mathi.Max(dx, dy)
}
