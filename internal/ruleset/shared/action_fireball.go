package shared

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/ruleset/base"
	"anvil/internal/tag"
)

type FireballAction struct {
	base.Action
}

func NewFireballAction(owner *core.Actor) FireballAction {
	tc := tag.NewContainer(tags.Spell, tags.Evocation)
	cost := map[tag.Tag]int{tags.SpellSlot3: 1, tags.Action: 1}
	a := FireballAction{
		Action: base.MakeAction(
			owner,
			"Fireball",
			tc,
			cost,
			30,
			4,
			[]core.DamageSource{{Times: 8, Sides: 6, Source: "Fireball", Tags: tag.NewContainer(tags.Fire)}},
		),
	}
	a.Tags().Add(tag.NewContainer(tags.Attack))
	return a
}

func (a FireballAction) Perform(pos []grid.Position, commitCost bool) {
	targets := a.targetsAt(pos[0])
	a.Owner().Log.Start(core.UseActionType, core.UseActionEvent{Action: a, Source: a.Owner(), Target: pos})
	a.Owner().Log.Add(core.TargetType, core.TargetEvent{Target: targets})
	defer a.Owner().Log.End()
	if commitCost {
		a.Commit()
	}
	dmg := a.Owner().DamageRoll(a.Damage(), false)
	for _, t := range targets {
		currDmg := dmg.Clone()
		save := t.SaveThrow(tags.Dexterity, a.Owner().SpellSaveDC())
		if save.Success {
			currDmg.HalveDamage(tags.Fire, "Saving throw")
		}
		t.TakeDamage(currDmg)
	}
}

func (a FireballAction) ValidPositions(from grid.Position) []grid.Position {
	if !a.CanAfford() {
		return []grid.Position{}
	}
	valid := []grid.Position{}
	shape := shapes.Sphere(from, a.CastRange())
	for _, pos := range shape {
		if !a.Owner().World.IsValidPosition(pos) {
			continue
		}
		if pos == from {
			continue
		}
		if !a.Owner().World.HasLineOfSight(from, pos) {
			continue
		}
		valid = append(valid, pos)
	}
	return valid
}

func (a FireballAction) targetsAt(pos grid.Position) []*core.Actor {
	valid := a.AffectedPositions([]grid.Position{pos})
	targets := make([]*core.Actor, 0)
	for _, p := range valid {
		cell := a.Owner().World.At(p)
		if cell == nil {
			continue
		}
		occupant := cell.Occupant()
		if occupant == nil {
			continue
		}

		if occupant.IsDead() {
			continue
		}
		targets = append(targets, occupant)
	}
	return targets
}

func (a FireballAction) AffectedPositions(tar []grid.Position) []grid.Position {
	valid := a.Owner().World.FloodFill(tar[0], a.Reach())
	return valid
}
