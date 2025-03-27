package pathfinding

import (
	"anvil/internal/grid"
	"slices"

	"github.com/adam-lavrik/go-imath/ix"
)

func (pf *Pathfinding) FindPath(start grid.Position, end grid.Position) []grid.Position {
	pf.reset()
	open := []*Node{}
	closed := []*Node{}
	startNode, _ := pf.grid.At(start)
	endNode, _ := pf.grid.At(end)
	open = append(open, startNode)
	startNode.SetGCost(0)
	startNode.SetHCost(pf.distance(start, end))
	for len(open) > 0 {
		current := pf.lowestFCost(open)
		if current == endNode {
			return pf.calculatePath(endNode)
		}
		open = slices.DeleteFunc(open, func(n *Node) bool {
			return n == current
		})
		closed = append(closed, current)
		for _, neighbour := range pf.neighbours(current) {
			if slices.Contains(closed, neighbour) {
				continue
			}
			tentativeGCost := current.GCost() + pf.distance(current.Position(), neighbour.Position())
			if tentativeGCost < neighbour.GCost() {
				neighbour.SetParent(current)
				neighbour.SetGCost(tentativeGCost)
				neighbour.SetHCost(pf.distance(neighbour.Position(), end))
				if !slices.Contains(open, neighbour) {
					open = append(open, neighbour)
				}
			}
		}
	}
	return nil
}

func (pf *Pathfinding) reset() {
	for x := 0; x < pf.width; x++ {
		for y := 0; y < pf.height; y++ {
			node, _ := pf.grid.At(grid.NewPosition(x, y))
			node.reset()
		}
	}
}

func (pf *Pathfinding) distance(a grid.Position, b grid.Position) int {
	xd := ix.Abs(a.X - b.X)
	yd := ix.Abs(a.Y - b.Y)
	remaining := ix.Abs(xd - yd)
	return MOVE_DIAGONAL_COST*ix.Min(xd, yd) + MOVE_STRAIGHT_COST*remaining
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
		pos := offset.Add(node.Position())
		neighbour, ok := pf.grid.At(pos)
		if !ok {
			continue
		}
		if !neighbour.IsWalkable() {
			continue
		}
		// For diagonal moves, check if both adjacent cells are walkable
		dx := pos.X - node.Position().X
		dy := pos.Y - node.Position().Y
		if dx != 0 && dy != 0 {
			horizontalPos := node.Position().Add(grid.Position{X: dx, Y: 0})
			verticalPos := node.Position().Add(grid.Position{X: 0, Y: dy})
			horizontalNode, ok := pf.grid.At(horizontalPos)
			if !ok || !horizontalNode.IsWalkable() {
				continue
			}
			verticalNode, ok := pf.grid.At(verticalPos)
			if !ok || !verticalNode.IsWalkable() {
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

func (pf *Pathfinding) calculatePath(end *Node) []grid.Position {
	result := Result{}
	current := end
	for current != nil {
		result.path = append(result.path, current.Position())
		current = current.Parent()
	}
	result.cost = end.FCost()
	slices.Reverse(result.path)
	return result.path
}
