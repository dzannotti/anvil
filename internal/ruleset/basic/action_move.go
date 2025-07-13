package basic

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type MoveAction struct {
	owner     *core.Actor
	archetype string
	id        string
	name      string
	tags      tag.Container
	cost      map[tag.Tag]int
	castRange int
	reach     int
}

func NewMoveAction(owner *core.Actor) *MoveAction {
	a := &MoveAction{
		owner:     owner,
		archetype: "",
		id:        "",
		name:      "Move",
		tags:      tag.ContainerFromTag(tags.Move),
		cost:      map[tag.Tag]int{tags.ResourceSpeed: 1},
		castRange: 0,
		reach:     0,
	}
	return a
}

func (a MoveAction) Owner() *core.Actor {
	return a.owner
}

func (a MoveAction) Archetype() string {
	return a.archetype
}

func (a MoveAction) ID() string {
	return a.id
}

func (a MoveAction) Name() string {
	return a.name
}

func (a MoveAction) Tags() *tag.Container {
	return &a.tags
}

func (a MoveAction) Cost() map[tag.Tag]int {
	return a.cost
}

func (a MoveAction) Reach() int {
	return a.reach
}

func (a MoveAction) CastRange() int {
	return a.castRange
}

func (a MoveAction) CanAfford() bool {
	return a.owner.Resources.CanAfford(a.cost)
}

func (a MoveAction) Commit() {
	if !a.CanAfford() {
		panic("Attempt to commit action without affording cost")
	}

	for tag, amount := range a.cost {
		a.owner.ConsumeResource(tag, amount)
	}
}

func (a MoveAction) AverageDamage() int {
	return 0
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
		src.ConsumeResource(tags.ResourceWalkSpeed, 1)
		src.Move(node, a)
	}
}

func (a MoveAction) AffectedPositions(tar []grid.Position) []grid.Position {
	return []grid.Position{a.Owner().Position, tar[0]}
}

func (a MoveAction) ValidPositions(from grid.Position) []grid.Position {
	speed := a.owner.Resources.Remaining(tags.ResourceWalkSpeed)
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
