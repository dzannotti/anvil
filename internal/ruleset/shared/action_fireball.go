package shared

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/ruleset/base"
	"anvil/internal/tag"

	"github.com/google/uuid"
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
			"fireball",
			uuid.New().String(),
			"Fireball",
			tc,
			cost,
			30,
			4,
		),
	}
	a.Tags().Add(tag.NewContainer(tags.Attack))
	return a
}

func (a FireballAction) Perform(pos []grid.Position) {
	targets := a.targetsAt(pos[0])
	a.Owner().Dispatcher.Begin(core.UseActionEvent{Action: a, Source: a.Owner(), Target: pos})
	a.Owner().Dispatcher.Emit(core.TargetEvent{Target: targets})
	defer a.Owner().Dispatcher.End()
	a.Commit()
	dmg := a.Owner().DamageRoll(a, false)
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
	shape := shapes.Circle(from, a.CastRange())
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

func (a FireballAction) Damage() *expression.Expression {
	expr := expression.FromDamageDice(8, 6, "Fireball", tag.NewContainer(tags.Fire))
	return &expr
}

func (a FireballAction) Tags() *tag.Container {
	fireTags := tag.NewContainer(tags.Fire)
	return &fireTags
}

func (a FireballAction) AverageDamage() int {
	return a.Damage().ExpectedValue()
}
