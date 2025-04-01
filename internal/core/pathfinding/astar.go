package pathfinding

import (
	"anvil/internal/grid"
	"slices"

	"github.com/adam-lavrik/go-imath/ix"
)

const (
	MoveDiagonalCost = 10
	MoveStraightCost = 10
)

func (pf *Pathfinding) FindPath(start grid.Position, end grid.Position) (*Result, bool) {
	pf.reset()
	open := []*Node{}
	closed := []*Node{}
	startNode, _ := pf.grid.At(start)
	endNode, _ := pf.grid.At(end)
	open = append(open, startNode)
	startNode.GCost = 0
	startNode.HCost = pf.distance(start, end)
	for len(open) > 0 {
		current := pf.lowestFCost(open)
		if current == endNode {
			return pf.calculatePath(endNode), true
		}
		open = slices.DeleteFunc(open, func(n *Node) bool {
			return n == current
		})
		closed = append(closed, current)
		for _, neighbour := range pf.neighbours(current) {
			if slices.Contains(closed, neighbour) {
				continue
			}
			tentativeGCost := current.GCost + pf.distance(current.Position, neighbour.Position)
			if tentativeGCost < neighbour.GCost {
				neighbour.Parent = current
				neighbour.GCost = tentativeGCost
				neighbour.HCost = pf.distance(neighbour.Position, end)
				if !slices.Contains(open, neighbour) {
					open = append(open, neighbour)
				}
			}
		}
	}
	return &Result{}, false
}

func (pf *Pathfinding) reset() {
	for x := 0; x < pf.width; x++ {
		for y := 0; y < pf.height; y++ {
			node, _ := pf.grid.At(grid.Position{X: x, Y: y})
			node.reset()
		}
	}
}

func (pf *Pathfinding) distance(a grid.Position, b grid.Position) int {
	xd := ix.Abs(a.X - b.X)
	yd := ix.Abs(a.Y - b.Y)
	remaining := ix.Abs(xd - yd)
	return MoveDiagonalCost*ix.Min(xd, yd) + MoveStraightCost*remaining
}

func (pf *Pathfinding) neighbours(node *Node) []*Node {
	offset := []grid.Position{
		{X: 0, Y: -1},  // up
		{X: 0, Y: 1},   // down
		{X: 1, Y: 0},   // right
		{X: -1, Y: 0},  // left
		{X: 1, Y: -1},  // up right
		{X: 1, Y: 1},   // down right
		{X: -1, Y: 1},  // down left
		{X: -1, Y: -1}, // up left
	}
	neighbours := make([]*Node, 0)
	for _, offset := range offset {
		pos := offset.Add(node.Position)
		neighbour, ok := pf.grid.At(pos)
		if !ok {
			continue
		}
		if !neighbour.Walkable {
			continue
		}
		// For diagonal moves, check if both adjacent cells are walkable
		dx := pos.X - node.Position.X
		dy := pos.Y - node.Position.Y
		if dx != 0 && dy != 0 {
			horizontalPos := node.Position.Add(grid.Position{X: dx, Y: 0})
			verticalPos := node.Position.Add(grid.Position{X: 0, Y: dy})
			horizontalNode, ok := pf.grid.At(horizontalPos)
			if !ok || !horizontalNode.Walkable {
				continue
			}
			verticalNode, ok := pf.grid.At(verticalPos)
			if !ok || !verticalNode.Walkable {
				continue
			}
		}
		neighbours = append(neighbours, neighbour)
	}
	return neighbours
}

func (pf *Pathfinding) lowestFCost(nodes []*Node) *Node {
	slices.SortFunc(nodes, func(a, b *Node) int {
		return a.FCost() - b.FCost()
	})
	return nodes[0]
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
