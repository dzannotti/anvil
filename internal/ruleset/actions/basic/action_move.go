package basic

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type MoveAction struct {
	Action
}

func NewMoveAction(owner *core.Actor) *MoveAction {
	a := &MoveAction{
		Action: Action{
			owner: owner,
			name:  "Move",
			cost:  map[tag.Tag]int{tags.Speed: 1},
			tags:  tag.NewContainer(tags.Move),
		},
	}
	return a
}

func (a MoveAction) Perform(pos []grid.Position) {
	src := a.owner
	world := src.World
	path, ok := world.FindPath(src.Position, pos[0])
	if !ok {
		panic("attempted to move to unreachable location - this should never happen")
	}

	src.Dispatcher.Begin(core.MoveEvent{World: world, Source: src, From: src.Position, To: pos[0], Path: path})
	defer src.Dispatcher.End()
	positions := path.Positions()
	for _, node := range positions[1:] {
		src.ConsumeResource(tags.WalkSpeed, 1)
		src.Move(node, a)
	}
}

func (a MoveAction) AffectedPositions(tar []grid.Position) []grid.Position {
	return []grid.Position{a.Owner().Position, tar[0]}
}

func (a MoveAction) ValidPositions(from grid.Position) []grid.Position {
	speed := a.owner.Resources.Remaining(tags.WalkSpeed)
	shape := shapes.Circle(from, speed)
	valid := make([]grid.Position, 0)
	for _, pos := range shape {
		if !a.owner.World.IsValidPosition(pos) {
			continue
		}

		if pos == from {
			continue
		}

		cell := a.owner.World.Grid.At(pos)
		if cell.IsOccupied() {
			continue
		}

		path, ok := a.owner.World.FindPath(from, pos)
		if !ok || path.Speed() > speed {
			continue
		}

		valid = append(valid, pos)
	}
	return valid
}
