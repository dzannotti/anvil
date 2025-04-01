package pathfinding

import (
	"anvil/internal/grid"
	"math"
)

type Node struct {
	Position grid.Position
	Parent   *Node
	GCost    int
	HCost    int
	Walkable bool
}

func (n Node) FCost() int {
	return n.GCost + n.HCost
}

func (n *Node) reset() {
	n.Parent = nil
	n.GCost = math.MaxInt
	n.HCost = 0
}
