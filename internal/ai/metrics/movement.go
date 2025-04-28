package metrics

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"math"
)

type Movement struct{}

func (d Movement) Evaluate(world *core.World, actor *core.Actor, action core.Action, pos grid.Position, affected []grid.Position) int {
	if !action.Tags().MatchTag(tags.Move) {
		return 0
	}

	if actor.Position == pos {
		return 0.0
	}

	lookAhead := 4
	speed := actor.Resources.Remaining(tags.WalkSpeed)
	enemies := world.ActorsInRange(pos, speed*lookAhead, func(other *core.Actor) bool { return other.IsHostileTo(actor) })

	score := 0

	distNow, distThen := d.closestAt(actor, pos, enemies)

	if distThen >= distNow {
		return 0
	}

	compression := float32(distNow-distThen) / (float32(distNow) + 0.001)
	distWeight := compression * 5
	targetCount := len(enemiesAffected(world, actor, action.ValidPositions(pos)))
	currentTargetCount := len(enemiesAffected(world, actor, action.ValidPositions(actor.Position)))

	if currentTargetCount > 0 && targetCount <= currentTargetCount {
		return 0
	}

	targetWeight := float32(targetCount) / float32(len(enemies))
	aooPenalty := float32(d.estimateOpportunityAttackDamageAt(pos))

	score = int(distWeight*4.0 + targetWeight*6.0 - aooPenalty*0.5 + 0.5)

	if score < 1 {
		return 0
	}

	return score
}

func (a Movement) estimateOpportunityAttackDamageAt(_ grid.Position) float64 {
	// TODO: Implement AOO here
	return 0.0
}

func (a Movement) closestAt(src *core.Actor, dst grid.Position, enemies []*core.Actor) (int, int) {
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
