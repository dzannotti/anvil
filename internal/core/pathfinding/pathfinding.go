package pathfinding

import (
	"container/heap"
	"math"
	"slices"

	"anvil/internal/grid"

	"github.com/adam-lavrik/go-imath/ix"
)

type Result struct {
	Path  []grid.Position
	Cost  int
	Speed int
}

type node struct {
	pos    grid.Position
	fScore int
	index  int // used by heap
}

type nodeHeap []*node

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].fScore < h[j].fScore }
func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *nodeHeap) Push(x any) {
	n := x.(*node)
	n.index = len(*h)
	*h = append(*h, n)
}
func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
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

func FindPath(start grid.Position, end grid.Position, width int, height int, movementCost func(grid.Position) int) (*Result, bool) {
	// Use flat arrays for better cache locality
	size := width * height
	gCost := make([]int, size)
	cameFrom := make([]*grid.Position, size)

	// Index conversion function
	idx := func(p grid.Position) int {
		return p.Y*width + p.X
	}

	// Initialize costs
	for i := range gCost {
		gCost[i] = math.MaxInt
	}

	inBounds := func(p grid.Position) bool {
		return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
	}

	open := &nodeHeap{}
	heap.Init(open)

	startNode := &node{
		pos:    start,
		fScore: heuristic(start, end),
	}
	heap.Push(open, startNode)
	gCost[idx(start)] = 0

	for open.Len() > 0 {
		current := heap.Pop(open).(*node)
		currentIdx := idx(current.pos)

		if current.pos == end {
			return reconstructPath(cameFrom, end, width, gCost[idx(end)]), true
		}

		// If we've found a better path to this node already, skip it
		if current.fScore > gCost[currentIdx]+heuristic(current.pos, end) {
			continue
		}

		for i, offset := range offsets {
			neighborPos := current.pos.Add(offset)
			if !inBounds(neighborPos) {
				continue
			}

			cost := movementCost(neighborPos)
			if cost == math.MaxInt {
				continue // impassable
			}

			// Check diagonal wall cutting only for diagonal moves
			isDiagonal := i >= 4 // First 4 are cardinal, last 4 are diagonal
			if isDiagonal {
				dx := offset.X
				dy := offset.Y
				adj1 := grid.Position{X: current.pos.X + dx, Y: current.pos.Y}
				adj2 := grid.Position{X: current.pos.X, Y: current.pos.Y + dy}
				if movementCost(adj1) == math.MaxInt || movementCost(adj2) == math.MaxInt {
					continue
				}
			}

			// Use the same cost factor for diagonals as for cardinal directions
			moveCost := cost * 10

			neighborIdx := idx(neighborPos)
			tentativeG := gCost[currentIdx] + moveCost

			if tentativeG < gCost[neighborIdx] {
				gCost[neighborIdx] = tentativeG
				cameFrom[neighborIdx] = &current.pos

				// Create a new node (no pooling to maintain thread safety)
				newNode := &node{
					pos:    neighborPos,
					fScore: tentativeG + heuristic(neighborPos, end),
				}
				heap.Push(open, newNode)
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
	dx := ix.Abs(a.X - b.X)
	dy := ix.Abs(a.Y - b.Y)
	// Chebyshev distance - maximum of dx and dy
	return 10 * ix.Max(dx, dy)
}
