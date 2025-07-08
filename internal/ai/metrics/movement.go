package metrics

import (
	"math"

	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type Movement struct{}

func (m Movement) Evaluate(
	world *core.World,
	actor *core.Actor,
	action core.Action,
	pos grid.Position,
	_ []grid.Position,
) int {
	if !action.Tags().MatchTag(tags.Move) {
		return 0
	}

	if actor.Position == pos {
		return 0
	}

	lookAhead := 4
	speed := actor.Resources.Remaining(tags.WalkSpeed)
	enemies := world.ActorsInRange(
		pos,
		speed*lookAhead,
		func(other *core.Actor) bool { return other.IsHostileTo(actor) },
	)

	distNow, distThen := m.closestAt(actor, pos, enemies)

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

	targetWeight := float32(targetCount) / (float32(len(enemies)) + 0.001)
	aooPenalty := float32(m.estimateOpportunityAttackDamageAt(pos))

	score := int(distWeight*4.0 + targetWeight*6.0 - aooPenalty*0.5 + 0.5)

	if score < 1 {
		return 0
	}

	return score
}

func (m Movement) estimateOpportunityAttackDamageAt(_ grid.Position) float64 {
	// TODO: Implement AOO here
	return 0
}

func (m Movement) closestAt(src *core.Actor, dst grid.Position, enemies []*core.Actor) (int, int) {
	world := src.World
	distNow := math.MaxInt
	distThen := math.MaxInt
	for _, enemy := range enemies {
		if path, ok := world.FindPath(src.Position, enemy.Position); ok && int(path.TotalCost) < distNow {
			distNow = int(path.TotalCost)
		}
		if path, ok := world.FindPath(dst, enemy.Position); ok && int(path.TotalCost) < distThen {
			distThen = int(path.TotalCost)
		}
	}
	return distNow, distThen
}
