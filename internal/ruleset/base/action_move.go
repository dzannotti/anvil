package base

import (
	"math"

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
			tags:  tag.ContainerFromTag(tags.Move),
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
	src.Log.Start(core.MoveType, core.MoveEvent{World: world, Source: src, From: src.Position, To: pos[0], Path: path})
	defer src.Log.End()
	for _, node := range path.Path[1:] {
		src.Resources.Consume(tags.WalkSpeed, 1)
		src.Log.Add(core.SpendResourceType, core.SpendResourceEvent{Source: src, Resource: tags.WalkSpeed, Amount: 1})
		src.Log.Start(core.MoveStepType, core.MoveStepEvent{World: world, Source: src, From: src.Position, To: node})
		// TODO: Implement AOO here
		world.RemoveOccupant(src.Position, src)
		src.Position = node
		src.World.AddOccupant(node, src)
		src.Log.End()
	}
}

func (a MoveAction) ScoreAt(dst grid.Position) float32 {
	src := a.owner
	world := src.World
	if src.Position == dst {
		return 0.0
	}
	lookAhead := 4
	speed := src.Resources.Remaining(tags.WalkSpeed)
	enemies := world.ActorsInRange(dst, speed*lookAhead, func(other *core.Actor) bool { return other.IsHostileTo(src) })

	score := float32(0)

	distNow, distThen := a.closestAt(dst, enemies)

	if distThen >= distNow {
		return 0.0
	}

	compression := float32(distNow-distThen) / (float32(distNow) + 0.001)
	distWeight := compression * 0.5
	score = score + distWeight*0.3
	targetCount := src.TargetCountAt(dst)
	currentTargetCount := src.TargetCountAt(src.Position)

	if currentTargetCount > 0 && targetCount <= currentTargetCount {
		return 0.0
	}

	targetWeight := float32(targetCount) / float32(len(enemies))

	score = score + targetWeight*0.6

	// AOO penalty
	score -= float32(a.estimateOpportunityAttackDamageAt(dst)) * 1.1

	if score < 0.01 {
		return 0.0
	}

	return score
}

func (a MoveAction) TargetCountAt(_ grid.Position) int {
	return 0
}

func (a MoveAction) AffectedPositions(tar []grid.Position) []grid.Position {
	return []grid.Position{a.Owner().Position, tar[0]}
}

func (a MoveAction) estimateOpportunityAttackDamageAt(_ grid.Position) float64 {
	// TODO: Implement AOO here
	return 0.0
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
		cell, _ := a.owner.World.Grid.At(pos)
		if cell.IsOccupied() {
			continue
		}
		path, ok := a.owner.World.FindPath(from, pos)
		if !ok || path.Speed > speed {
			continue
		}
		valid = append(valid, pos)
	}
	return valid
}

func (a MoveAction) closestAt(dst grid.Position, enemies []*core.Actor) (int, int) {
	src := a.owner
	world := src.World
	distNow := math.MaxInt
	distThen := math.MaxInt
	for _, enemy := range enemies {
		if path, ok := world.FindPath(src.Position, enemy.Position); ok && path.Speed < distNow {
			distNow = path.Speed
		}
		if path, ok := world.FindPath(dst, enemy.Position); ok && path.Speed < distThen {
			distThen = path.Speed
		}
	}
	return distNow, distThen
}
