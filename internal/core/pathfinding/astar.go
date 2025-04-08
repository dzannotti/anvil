package pathfinding

import (
	"container/heap"
	"slices"

	"github.com/adam-lavrik/go-imath/ix"

	"anvil/internal/grid"
)

const (
	MoveDiagonalCost = 1
	MoveStraightCost = 1
)

var pathfindingOffsets = []grid.Position{
	{X: 0, Y: -1},  // up
	{X: 0, Y: 1},   // down
	{X: 1, Y: 0},   // right
	{X: -1, Y: 0},  // left
	{X: 1, Y: -1},  // up right
	{X: 1, Y: 1},   // down right
	{X: -1, Y: 1},  // down left
	{X: -1, Y: -1}, // up left
}

func (pf *Pathfinding) FindPath(start, end grid.Position) (*Result, bool) {
	pf.reset()
	startNode, _ := pf.grid.At(start)
	endNode, _ := pf.grid.At(end)

	openPQ := make(PriorityQueue, 0, 1024)
	heap.Init(&openPQ)
	heap.Push(&openPQ, startNode)

	openSet := map[*Node]struct{}{startNode: {}}
	closedSet := map[*Node]struct{}{}

	startNode.GCost = 0
	startNode.HCost = pf.distance(start, end)

	for openPQ.Len() > 0 {
		current := heap.Pop(&openPQ).(*Node)
		delete(openSet, current)
		closedSet[current] = struct{}{}

		if current == endNode {
			return pf.calculatePath(endNode), true
		}

		for _, neighbour := range pf.neighbours(current) {
			if _, closed := closedSet[neighbour]; closed {
				continue
			}

			tentativeG := current.GCost + pf.distance(current.Position, neighbour.Position)
			if tentativeG < neighbour.GCost {
				neighbour.Parent = current
				neighbour.GCost = tentativeG
				neighbour.HCost = pf.distance(neighbour.Position, end)

				if _, open := openSet[neighbour]; !open {
					heap.Push(&openPQ, neighbour)
					openSet[neighbour] = struct{}{}
				}
			}
		}
	}

	return nil, false
}

func (pf *Pathfinding) reset() {
	for x := 0; x < pf.width; x++ {
		for y := 0; y < pf.height; y++ {
			node, _ := pf.grid.At(grid.Position{X: x, Y: y})
			node.reset()
		}
	}
}

func (pf *Pathfinding) distance(a, b grid.Position) int {
	xd := ix.Abs(a.X - b.X)
	yd := ix.Abs(a.Y - b.Y)
	remaining := ix.Abs(xd - yd)
	return MoveDiagonalCost*ix.Min(xd, yd) + MoveStraightCost*remaining
}

func (pf *Pathfinding) neighbours(node *Node) []*Node {
	neighbours := make([]*Node, 0, len(pathfindingOffsets))
	for _, offset := range pathfindingOffsets {
		pos := offset.Add(node.Position)
		neighbour, ok := pf.grid.At(pos)
		if !ok || !neighbour.Walkable {
			continue
		}

		// Diagonal movement check
		dx := pos.X - node.Position.X
		dy := pos.Y - node.Position.Y
		if dx != 0 && dy != 0 {
			horizontal, okH := pf.grid.At(node.Position.Add(grid.Position{X: dx, Y: 0}))
			vertical, okV := pf.grid.At(node.Position.Add(grid.Position{X: 0, Y: dy}))
			if !okH || !okV || !horizontal.Walkable || !vertical.Walkable {
				continue
			}
		}

		neighbours = append(neighbours, neighbour)
	}
	return neighbours
}

func (pf *Pathfinding) calculatePath(end *Node) *Result {
	result := Result{}
	current := end
	for current != nil {
		result.Path = append(result.Path, current.Position)
		current = current.Parent
	}
	result.Cost = end.FCost()
	slices.Reverse(result.Path)
	return &result
}
