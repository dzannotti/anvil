package base

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

func NewMoveAction(owner *core.Actor) MoveAction {
	a := MoveAction{
		Action: Action{
			owner: owner,
			name:  "Move",
			cost:  map[tag.Tag]int{tags.Speed: 1},
			tags:  tag.ContainerFromTag(tags.Move),
		},
	}
	return a
}

func (a MoveAction) Perform(pos []grid.Position) {
	src := a.owner
	world := src.World
	path, ok := world.Navigation.FindPath(src.Position, pos[0])
	if !ok {
		panic("attempted to move to unreachable location - this should never happen")
	}
	src.Log.Start(core.MoveEventType, core.MoveEvent{World: world, Source: src, From: src.Position, To: pos[0], Path: path})
	defer src.Log.End()
	for _, node := range path.Path[1:] {
		src.Resources.Consume(tags.Speed, 1)
		src.Log.Add(core.SpendResourceType, core.SpendResourceEvent{Source: src, Resource: tags.Speed, Amount: 1})
		src.Log.Start(core.MoveStepType, core.MoveStepEvent{World: world, Source: src, From: src.Position, To: node})
		defer src.Log.End()
		// TODO: Implement AOO here
		world.RemoveOccupant(src.Position, src)
		src.Position = node
		src.World.AddOccupant(node, src)
	}
}

func (a MoveAction) ScoreAt(dest grid.Position) *core.ScoredAction {
	src := a.owner
	if src.Position == dest {
		return nil
	}

	// depth = recursion in evaluation
	notMove := func(action core.Action, depth int) bool { return depth < 3 }
	bestAtDest := src.BestScoredActionAtWhere(dest, notMove, 0)
	bestAtCurrent := src.BestScoredActionAtWhere(actor.Position, notMove, 0)

	if bestAfter == nil {
		return nil
	}

	return nil
}

func (a MoveAction) estimateOpportunityAttackDamageAt(dst grid.Position) float64 {
	// TODO: Implement AOO here
	return 0.0
}

func (a MoveAction) ValidPositions(from grid.Position) []grid.Position {
	speed := a.owner.Resources.Remaining(tags.Speed)
	shape := shapes.Circle(from, speed)
	valid := make([]grid.Position, 0)
	for _, pos := range shape {
		if !a.owner.World.IsValidPosition(pos) {
			continue
		}
		if pos == from {
			continue
		}
		cell, _ := a.owner.World.Grid.At(pos)
		if cell.IsOccupied() {
			continue
		}
		path, ok := a.owner.World.Navigation.FindPath(from, pos)
		if !ok || path.Cost > speed {
			continue
		}
		valid = append(valid, pos)
	}
	return valid
}
