package shared

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/ruleset/base"
	"anvil/internal/tag"
)

var hitFriendliesPenalty = float32(-3)

type FireballAction struct {
	base.Action
	damage []core.DamageSource
}

func NewFireballAction(owner *core.Actor) FireballAction {
	tc := tag.ContainerFromTag(tags.Spell, tags.Evocation)
	cost := map[tag.Tag]int{tags.SpellSlot3: 1, tags.Action: 1}
	a := FireballAction{
		Action: base.MakeAction(owner, "Fireball", tc, cost),
		damage: []core.DamageSource{{Times: 8, Sides: 6, Source: "Fireball", Tags: tag.ContainerFromTag(tags.Fire)}},
	}
	a.Tags().Add(tag.ContainerFromTag(tags.Attack))
	return a
}

func (a FireballAction) Reach() int {
	return 4
}

func (a FireballAction) Range() int {
	return 30
}

func (a FireballAction) ScoreAt(pos grid.Position) float32 {
	targets := a.targetsAt(pos)
	if len(targets) == 0 {
		return 0.0
	}
	avgDmg := a.AverageDamage(a.damage)
	score := float32(0)
	for _, t := range targets {
		if t.IsDead() {
			continue
		}
		damageRatio := float32(avgDmg) / float32(t.HitPoints)
		if damageRatio > 1.0 {
			damageRatio = 1.0
		}
		lowHPPriority := (1 - t.HitPointsNormalized()) * 0.5
		if !a.Owner().IsHostileTo(t) {
			damageRatio = damageRatio * hitFriendliesPenalty
			lowHPPriority = lowHPPriority * hitFriendliesPenalty
		}
		score = score + damageRatio + lowHPPriority
	}
	if score > 1 {
		score = 1.0
	}
	return score
}

func (a FireballAction) Perform(pos []grid.Position) {
	targets := a.targetsAt(pos[0])
	a.Owner().Log.Start(core.UseActionType, core.UseActionEvent{Action: a, Source: a.Owner(), Target: pos})
	a.Owner().Log.Add(core.TargetType, core.TargetEvent{Target: targets})
	defer a.Owner().Log.End()
	a.Commit()
	dmg := a.Owner().DamageRoll(a.damage, false)
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
	shape := shapes.Sphere(from, a.Range())
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

func (a FireballAction) TargetCountAt(pos grid.Position) int {
	targets := a.targetsAt(pos)
	count := 0
	for _, t := range targets {
		if !a.Owner().IsHostileTo(t) {
			count = count + int(hitFriendliesPenalty)
			continue
		}
		if !a.Owner().World.HasLineOfSight(a.Owner().Position, t.Position) {
			continue
		}
		count = count + 1
	}
	return max(count, 0)
}

func (a FireballAction) targetsAt(pos grid.Position) []*core.Actor {
	valid := a.AffectedPositions([]grid.Position{pos})
	targets := make([]*core.Actor, 0)
	for _, p := range valid {
		cell, ok := a.Owner().World.At(p)
		if !ok {
			continue
		}
		occupant, ok := cell.Occupant()
		if !ok {
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
