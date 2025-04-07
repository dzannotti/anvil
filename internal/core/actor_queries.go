package core

import (
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"

	"github.com/adam-lavrik/go-imath/ix"
)

func (a Actor) Enemies() []*Actor {
	opponents := TeamEnemies
	if a.Team == TeamEnemies {
		opponents = TeamPlayers
	}
	enemies := make([]*Actor, 0, len(a.Encounter.Actors))
	for _, c := range a.Encounter.Actors {
		if opponents == c.Team {
			enemies = append(enemies, c)
		}
	}
	return enemies
}

func (a Actor) HitPointsNormalized() float32 {
	return float32(a.HitPoints) / float32(a.MaxHitPoints)
}

func (a Actor) HasCondition(t tag.Tag, src *Effect) bool {
	return a.Conditions.Has(t, src)
}

func (a Actor) IsDead() bool {
	return a.HasCondition(tags.Dead, nil)
}

func (a Actor) CanAct() bool {
	return !a.Conditions.Match(tags.Incapacitated)
}

func (a Actor) TargetCountAt(pos grid.Position) int {
	c := 0
	for _, a := range a.Actions {
		c = ix.Max(a.TargetCountAt(pos), c)
	}
	return c
}
