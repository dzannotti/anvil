package base

import (
	"slices"

	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type AttackAction struct {
	Action
	reach        int
	DamageSource []core.DamageSource
}

func NewAttackAction(owner *core.Actor, name string, ds []core.DamageSource, reach int, t ...tag.Tag) AttackAction {
	a := AttackAction{
		Action: Action{
			owner: owner,
			name:  name,
			cost:  map[tag.Tag]int{tags.Action: 1},
			tags:  tag.ContainerFromTag(t...),
		},
		reach:        reach,
		DamageSource: ds,
	}
	a.tags.Add(tag.ContainerFromTag(tags.Melee, tags.Attack))
	a.WithScorer(a.Score)
	return a
}

func (a AttackAction) Perform(pos []grid.Position) {
	a.Commit()
	target, _ := a.owner.World.ActorAt(pos[0])
	a.owner.Log.Start(core.UseActionType, core.UseActionEvent{Action: a, Source: *a.owner, Target: *target})
	defer a.owner.Log.End()
	result := a.owner.AttackRoll(target, tag.Container{})
	if result.Success {
		dmg := a.owner.DamageRoll(a.DamageSource, result.Critical)
		target.TakeDamage(dmg.Value)
	}
}

func (a AttackAction) Score(pos grid.Position) float32 {
	target, _ := a.owner.World.ActorAt(pos)
	if target == nil {
		return 0
	}
	return 0.5 + (1-target.HitPointsNormalized())*0.5
}

func (a AttackAction) ValidPositions(from grid.Position) []grid.Position {
	if !a.CanAfford() {
		return []grid.Position{}
	}
	reach := a.reach
	shape := shapes.Sphere(from, reach)
	valid := make([]grid.Position, 0)
	enemies := a.owner.Enemies()
	for _, pos := range shape {
		if !a.owner.World.IsValidPosition(pos) {
			continue
		}
		if pos == from {
			continue
		}
		other, ok := a.owner.World.ActorAt(pos)
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
