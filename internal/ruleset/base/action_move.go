package base

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
	"math"
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
	src.Log.Start(core.MoveType, core.MoveEvent{World: world, Source: src, From: src.Position, To: pos[0], Path: path})
	defer src.Log.End()
	for _, node := range path.Path[1:] {
		src.Resources.Consume(tags.Speed, 1)
		src.Log.Add(core.SpendResourceType, core.SpendResourceEvent{Source: src, Resource: tags.Speed, Amount: 1})
		src.Log.Start(core.MoveStepType, core.MoveStepEvent{World: world, Source: src, From: src.Position, To: node})
		// TODO: Implement AOO here
		world.RemoveOccupant(src.Position, src)
		src.Position = node
		src.World.AddOccupant(node, src)
		src.Log.End()
	}
}

func (a MoveAction) ScoreAt(dest grid.Position) *core.ScoredAction {
	src := a.owner
	world := src.World
	if src.Position == dest {
		return nil
	}
	speed := src.Resources.Remaining(tags.WalkSpeed)
	lookAhead := 4
	enemies := world.ActorsInRange(dest, speed*lookAhead, func(other *core.Actor) bool { return other.Team != src.Team })

	distNow := math.MaxInt
	distThen := math.MaxInt
	for _, enemy := range enemies {
		if path, ok := world.Navigation.FindPath(src.Position, enemy.Position); ok && path.Cost < distNow {
			distNow = path.Cost
		}
		if path, ok := world.Navigation.FindPath(dest, enemy.Position); ok && path.Cost < distThen {
			distThen = path.Cost
		}
	}

	if distThen >= distNow {
		return nil
	}

	targetCount := src.TargetCountAt(dest)
	currentTargetCount := src.TargetCountAt(src.Position)

	if currentTargetCount > 0 && targetCount <= currentTargetCount {
		return nil
	}

	compression := float32(distNow-distThen) / float32(distNow)
	distWeight := compression * 0.5
	targetWeight := float32(targetCount) / float32(len(enemies))

	score := distWeight*0.3 + targetWeight*0.6

	// Avoid movement spam
	score -= 0.05

	// AOO penalty
	score -= float32(a.estimateOpportunityAttackDamageAt(dest)) * 1.1

	if score < 0.01 {
		return nil
	}

	return &core.ScoredAction{
		Action:   a,
		Position: []grid.Position{dest},
		Score:    score,
	}
}

func (a MoveAction) TargetCountAt(pos grid.Position) int {
	return 0
}

func (a MoveAction) estimateOpportunityAttackDamageAt(dst grid.Position) float64 {
	// TODO: Implement AOO here
	return 0.0
}

func (a MoveAction) closestAt(pos grid.Position, enemies []*core.Actor) int {
	min := math.MaxInt
	for _, e := range enemies {
		path, ok := a.owner.World.Navigation.FindPath(pos, e.Position)
		if !ok {
			continue
		}
		if path.Cost < min {
			min = path.Cost
		}
	}
	return min
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
		path, ok := a.owner.World.Navigation.FindPath(from, pos)
		if !ok || path.Cost > speed {
			continue
		}
		valid = append(valid, pos)
	}
	return valid
}
