package pathfinding

import (
	"anvil/internal/grid"
	"math"
)

type Node struct {
	position grid.Position
	parent   *Node
	gCost    int
	hCost    int
	walkable bool
}

func NewNode(position grid.Position) Node {
	return Node{
		position: position,
		parent:   nil,
		gCost:    0,
		hCost:    0,
		walkable: true,
	}
}

func (n Node) Position() grid.Position {
	return n.position
}

func (n Node) GCost() int {
	return n.gCost
}

func (n Node) HCost() int {
	return n.hCost
}

func (n Node) FCost() int {
	return n.gCost + n.hCost
}

func (n Node) Parent() *Node {
	return n.parent
}

func (n *Node) SetParent(parent *Node) {
	n.parent = parent
}

func (n *Node) SetGCost(gCost int) {
	n.gCost = gCost
}

func (n *Node) SetHCost(hCost int) {
	n.hCost = hCost
}

func (n *Node) SetWalkable(walkable bool) {
	n.walkable = walkable
}

func (n Node) IsWalkable() bool {
	return n.walkable
}

func (n *Node) reset() {
	n.parent = nil
	n.gCost = math.MaxInt
	n.hCost = 0
}
