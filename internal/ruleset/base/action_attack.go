package base

import (
	"slices"

	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type ScoringFunc = func(pos grid.Position) float32

type AttackAction struct {
	Owner        *core.Actor
	name         string
	scorer       ScoringFunc
	DamageSource []core.DamageSource
}

func NewAttackAction(owner *core.Actor, name string, ds []core.DamageSource) AttackAction {
	a := AttackAction{
		Owner:        owner,
		name:         name,
		DamageSource: ds,
	}
	a.scorer = a.Score
	return a
}

func (a *AttackAction) WithScorer(s ScoringFunc) {
	a.scorer = s
}

func (a AttackAction) Name() string {
	return a.name
}

func (a AttackAction) AIAction(pos grid.Position) *core.AIAction {
	return &core.AIAction{
		Action:   a,
		Position: []grid.Position{pos},
		Score:    a.scorer(pos),
	}
}

func (a AttackAction) Perform(pos []grid.Position) {
	target, _ := a.Owner.World.ActorAt(pos[0])
	a.Owner.Log.Start(core.UseActionType, core.UseActionEvent{Action: a, Source: *a.Owner, Target: *target})
	defer a.Owner.Log.End()
	result := a.Owner.AttackRoll(target, tag.Container{})
	if result.Success {
		dmg := a.Owner.DamageRoll(a.DamageSource, result.Critical)
		target.TakeDamage(dmg.Value)
	}
}

func (a AttackAction) Score(pos grid.Position) float32 {
	target, _ := a.Owner.World.ActorAt(pos)
	if target == nil {
		return 0
	}
	return 0.5 + (1-target.HitPointsNormalized())*0.5
}

func (a AttackAction) ValidPositions(from grid.Position) []grid.Position {
	reach := 10
	shape := shapes.Sphere(from, reach)
	valid := make([]grid.Position, 0)
	enemies := a.Owner.Enemies()
	for _, pos := range shape {
		if !a.Owner.World.IsValidPosition(pos) {
			continue
		}
		if pos == from {
			continue
		}
		other, ok := a.Owner.World.ActorAt(pos)
		if !ok {
			continue
		}
		if !slices.Contains(enemies, other) {
			continue
		}
		if other.IsDead() {
			continue
		}
		valid = append(valid, pos)
	}
	return valid
}
