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
	reach  int
	damage []core.DamageSource
}

func NewAttackAction(owner *core.Actor, name string, ds []core.DamageSource, reach int, t ...tag.Tag) AttackAction {
	a := AttackAction{
		Action: Action{
			owner: owner,
			name:  name,
			cost:  map[tag.Tag]int{tags.Action: 1},
			tags:  tag.ContainerFromTag(t...),
		},
		reach:  reach,
		damage: ds,
	}
	a.tags.Add(tag.ContainerFromTag(tags.Attack))
	return a
}

func (a AttackAction) Perform(pos []grid.Position) {
	target, _ := a.owner.World.ActorAt(pos[0])
	a.owner.Log.Start(core.UseActionType, core.UseActionEvent{Action: a, Source: *a.owner, Target: *target})
	defer a.owner.Log.End()
	a.Commit()
	result := a.owner.AttackRoll(target, tag.Container{})
	if result.Success {
		dmg := a.owner.DamageRoll(a.damage, result.Critical)
		target.TakeDamage(*dmg)
	}
}

func (a AttackAction) ScoreAt(pos grid.Position) *core.ScoredAction {
	target, _ := a.owner.World.ActorAt(pos)
	if target == nil {
		return nil
	}
	avgDmg := a.AverageDamage(a.damage)
	damageRatio := float32(avgDmg) / float32(target.HitPoints)
	if damageRatio > 1.0 {
		damageRatio = 1.0
	}
	lowHPPriority := (1 - target.HitPointsNormalized()) * 0.5
	score := damageRatio + lowHPPriority
	return &core.ScoredAction{
		Action:   &a,
		Position: []grid.Position{pos},
		Score:    score,
	}
}

func (a AttackAction) ValidPositions(from grid.Position) []grid.Position {
	if !a.CanAfford() {
		return []grid.Position{}
	}
	shape := shapes.Sphere(from, a.reach)
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
