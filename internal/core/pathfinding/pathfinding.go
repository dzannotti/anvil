package pathfinding

import (
	"container/heap"
	"math"
	"slices"

	"anvil/internal/grid"

	"github.com/adam-lavrik/go-imath/ix"
)

type Result struct {
	Path []grid.Position
	Cost int
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
	inBounds := func(p grid.Position) bool {
		return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
	}

	gCost := make([][]int, width)
	cameFrom := make([][]*grid.Position, width)
	for x := range gCost {
		gCost[x] = make([]int, height)
		cameFrom[x] = make([]*grid.Position, height)
		for y := range gCost[x] {
			gCost[x][y] = math.MaxInt
		}
	}

	open := &nodeHeap{}
	heap.Init(open)

	startNode := &node{
		pos:    start,
		fScore: heuristic(start, end),
	}
	heap.Push(open, startNode)
	gCost[start.X][start.Y] = 0

	for open.Len() > 0 {
		current := heap.Pop(open).(*node)
		if current.pos == end {
			return reconstructPath(cameFrom, end, gCost[end.X][end.Y]), true
		}

		for _, offset := range offsets {
			neighborPos := current.pos.Add(offset)
			if !inBounds(neighborPos) {
				continue
			}

			cost := movementCost(neighborPos)
			if cost == math.MaxInt {
				continue // impassable
			}

			// prevent diagonal cuts through walls
			dx := neighborPos.X - current.pos.X
			dy := neighborPos.Y - current.pos.Y
			if dx != 0 && dy != 0 {
				adj1 := grid.Position{X: current.pos.X + dx, Y: current.pos.Y}
				adj2 := grid.Position{X: current.pos.X, Y: current.pos.Y + dy}
				if movementCost(adj1) == math.MaxInt || movementCost(adj2) == math.MaxInt {
					continue
				}
			}

			moveCost := 10 * cost
			tentativeG := gCost[current.pos.X][current.pos.Y] + moveCost
			if tentativeG < gCost[neighborPos.X][neighborPos.Y] {
				gCost[neighborPos.X][neighborPos.Y] = tentativeG
				cameFrom[neighborPos.X][neighborPos.Y] = &current.pos
				heap.Push(open, &node{
					pos:    neighborPos,
					fScore: tentativeG + heuristic(neighborPos, end),
				})
			}
		}
	}

	return nil, false
}

func heuristic(a, b grid.Position) int {
	xd := ix.Abs(a.X - b.X)
	yd := ix.Abs(a.Y - b.Y)
	return 10 * (xd + yd)
}

func reconstructPath(cameFrom [][]*grid.Position, end grid.Position, cost int) *Result {
	path := []grid.Position{end}
	curr := end
	for {
		prev := cameFrom[curr.X][curr.Y]
		if prev == nil {
			break
		}
		curr = *prev
		path = append(path, curr)
	}
	slices.Reverse(path)
	return &Result{
		Path: path,
		Cost: cost,
	}
}
