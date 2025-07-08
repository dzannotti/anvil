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
	reach int
}

func NewAttackAction(
	owner *core.Actor,
	name string,
	ds []core.DamageSource,
	reach int,
	tc tag.Container,
) *AttackAction {
	a := &AttackAction{
		Action: Action{
			owner:  owner,
			name:   name,
			cost:   map[tag.Tag]int{tags.Action: 1},
			tags:   tc,
			damage: ds,
		},
		reach: reach,
	}
	a.tags.Add(tag.NewContainer(tags.Attack))
	return a
}

func (a AttackAction) Perform(pos []grid.Position, commitCost bool) {
	target := a.owner.World.ActorAt(pos[0])
	a.owner.Dispatcher.Begin(core.UseActionType, core.UseActionEvent{Action: a, Source: a.owner, Target: pos})
	a.owner.Dispatcher.Emit(core.TargetType, core.TargetEvent{Target: []*core.Actor{target}})
	defer a.owner.Dispatcher.End()
	if commitCost {
		a.Commit()
	}
	result := a.owner.AttackRoll(target, a.tags)
	if result.Success {
		dmg := a.owner.DamageRoll(a.damage, result.Critical)
		target.TakeDamage(*dmg)
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
		other := a.owner.World.ActorAt(pos)
		if other == nil {
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

func (a AttackAction) AffectedPositions(tar []grid.Position) []grid.Position {
	return []grid.Position{tar[0]}
}
