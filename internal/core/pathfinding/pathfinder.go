package pathfinding

import (
	"math"
	"slices"

	"anvil/internal/grid"
	"anvil/internal/mathi"
)

type pathfinder struct {
	width, height int
	start, end    grid.Position
	movementCost  func(grid.Position) int
	gCost         []float64
	cameFrom      []*grid.Position
	open          *minHeap
}

func (pf *pathfinder) findPath() *Result {
	pf.initialize()

	for !pf.open.Empty() {
		current := pf.open.Pop()

		if current.pos == pf.end {
			return pf.reconstructPath()
		}

		pf.exploreNeighbors(current)
	}

	return &Result{Found: false}
}

func (pf *pathfinder) initialize() {
	for i := range pf.gCost {
		pf.gCost[i] = math.Inf(1)
	}

	startIdx := pf.posToIndex(pf.start)
	pf.gCost[startIdx] = 0
	pf.open.Push(&node{
		pos:    pf.start,
		fScore: float64(pf.heuristic(pf.start)),
	})
}

func (pf *pathfinder) exploreNeighbors(current *node) {
	for _, neighbor := range pf.getValidNeighbors(current.pos) {
		moveCost := pf.calculateMoveCost(current.pos, neighbor)

		neighborIdx := pf.posToIndex(neighbor)
		tentativeG := pf.gCost[pf.posToIndex(current.pos)] + moveCost

		if tentativeG < pf.gCost[neighborIdx] {
			pf.updateNeighbor(current.pos, neighbor, tentativeG)
		}
	}
}

func (pf *pathfinder) getValidNeighbors(pos grid.Position) []grid.Position {
	neighbors := make([]grid.Position, 0, 8)

	offsets := []grid.Position{
		{X: 0, Y: -1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 1, Y: 0}, // orthogonal
		{X: -1, Y: -1}, {X: 1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: 1}, // diagonal
	}

	for i, offset := range offsets {
		neighbor := pos.Add(offset)

		if !pf.inBounds(neighbor) {
			continue
		}

		if pf.movementCost(neighbor) == math.MaxInt {
			continue
		}

		// Prevent diagonal wall cutting
		if i >= 4 && pf.isDiagonalBlocked(pos, offset) {
			continue
		}

		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

func (pf *pathfinder) isDiagonalBlocked(pos, offset grid.Position) bool {
	adj1 := grid.Position{X: pos.X + offset.X, Y: pos.Y}
	adj2 := grid.Position{X: pos.X, Y: pos.Y + offset.Y}
	return pf.movementCost(adj1) == math.MaxInt || pf.movementCost(adj2) == math.MaxInt
}

func (pf *pathfinder) calculateMoveCost(from, to grid.Position) float64 {
	baseCost := float64(pf.movementCost(to))

	dx := mathi.Abs(to.X - from.X)
	dy := mathi.Abs(to.Y - from.Y)

	if dx > 0 && dy > 0 {
		return baseCost * DiagonalCost
	}
	return baseCost * NormalCost
}

func (pf *pathfinder) updateNeighbor(fromPos, toPos grid.Position, tentativeG float64) {
	toIdx := pf.posToIndex(toPos)
	pf.gCost[toIdx] = tentativeG
	pf.cameFrom[toIdx] = &fromPos

	pf.open.Push(&node{
		pos:    toPos,
		fScore: tentativeG + float64(pf.heuristic(toPos)),
	})
}

func (pf *pathfinder) reconstructPath() *Result {
	// Build path backwards first to get length
	positions := []grid.Position{pf.end}
	curr := pf.end

	for {
		prevPtr := pf.cameFrom[pf.posToIndex(curr)]
		if prevPtr == nil {
			break
		}
		curr = *prevPtr
		positions = append(positions, curr)
	}

	slices.Reverse(positions)

	// Build steps with metadata
	steps := make([]PathStep, len(positions))
	for i, pos := range positions {
		stepCost := 0.0
		if i > 0 {
			stepCost = pf.calculateMoveCost(positions[i-1], pos)
		}

		steps[i] = PathStep{
			Position: pos,
			GCost:    pf.gCost[pf.posToIndex(pos)],
			StepCost: stepCost,
			Distance: i,
		}
	}

	return &Result{
		Steps:     steps,
		TotalCost: pf.gCost[pf.posToIndex(pf.end)],
		Found:     true,
	}
}

func (pf *pathfinder) posToIndex(pos grid.Position) int {
	return pos.Y*pf.width + pos.X
}

func (pf *pathfinder) inBounds(pos grid.Position) bool {
	return pos.X >= 0 && pos.X < pf.width && pos.Y >= 0 && pos.Y < pf.height
}

func (pf *pathfinder) heuristic(pos grid.Position) int {
	dx := mathi.Abs(pos.X - pf.end.X)
	dy := mathi.Abs(pos.Y - pf.end.Y)
	return int(float64(mathi.Max(dx, dy)) * NormalCost)
}
