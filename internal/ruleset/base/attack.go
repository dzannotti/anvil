package base

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/grid"
	"anvil/internal/tag"
	"slices"
)

type AttackAction struct {
	Owner *core.Actor
}

func NewAttackAction(owner *core.Actor) AttackAction {
	return AttackAction{
		Owner: owner,
	}
}

func (a AttackAction) Name() string {
	return "Attack"
}

func (a AttackAction) Perform(pos []grid.Position) {
	target, _ := a.Owner.World.ActorAt(pos[0])
	a.Owner.Log.Start(core.UseActionEventType, core.UseActionEvent{Action: a, Source: *a.Owner, Target: *target})
	defer a.Owner.Log.End()
	result := a.Owner.AttackRoll(target, tag.Container{})
	if result.Success {
		target.TakeDamage(5)
	}
}

func (a AttackAction) AIAction(pos grid.Position) *core.AIAction {
	target, _ := a.Owner.World.ActorAt(pos)
	return &core.AIAction{
		Action:   a,
		Position: []grid.Position{pos},
		Score:    0.5 + (1-target.HitPointsNormalized())*0.5,
	}
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
