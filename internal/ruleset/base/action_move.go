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
	src.Log.Start(core.MoveEventType, core.MoveEvent{World: world, Source: src, From: src.Position, To: pos[0], Path: path})
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
	if src.Position == dest {
		return nil
	}
	speed := a.owner.Resources.Remaining(tags.WalkSpeed)
	lookAheadMoves := 4

	enemiesInRange := src.World.ActorsInRange(dest, speed*lookAheadMoves, func(other *core.Actor) bool { return other.Team != src.Team })

	score := float32(len(enemiesInRange)) * 0.3

	// reward getting closer to targets
	distNow := a.closestAt(src.Position, enemiesInRange)
	distThen := a.closestAt(dest, enemiesInRange)

	if distThen < distNow {
		// reward the improvement (relative change)
		compression := float32(distNow-distThen) / float32(distNow)
		score += compression * 0.4
	}

	if distThen > 0 {
		// the closer you are, the higher the bonus
		score += (1 / float32(distThen)) * 0.6
	}

	// avoid movement spam
	score -= 0.05

	// aoo penalty
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
