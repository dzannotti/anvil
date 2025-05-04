package core

import (
	"github.com/google/uuid"

	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func (a Actor) Enemies() []*Actor {
	enemies := make([]*Actor, 0, len(a.Encounter.Actors))
	for _, c := range a.Encounter.Actors {
		if a.IsHostileTo(c) {
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

func (a Actor) MatchCondition(t tag.Tag) bool {
	return a.Conditions.Match(t)
}

func (a Actor) IsDead() bool {
	return a.HasCondition(tags.Dead, nil)
}

func (a Actor) CanAct() bool {
	return !a.MatchCondition(tags.Incapacitated)
}

func (a Actor) SpellSaveDC() int {
	return 8 + a.Proficiencies.Bonus + stats.AttributeModifier(a.Attribute(a.SpellCastingSource).Value)
}

func (a Actor) IsHostileTo(o *Actor) bool {
	return a.Team != o.Team
}

func (a *Actor) ID() string {
	if a.id == "" {
		a.id = uuid.New().String()
	}
	return a.id
}

func (a Actor) HasAction(aa Action) bool {
	for _, ca := range a.Actions {
		if ca.Name() == aa.Name() {
			return true
		}
	}
	return false
}

func (a Actor) BestWeaponAttack() Action {
	var best Action
	bestDamage := 0
	for _, act := range a.Actions {
		if !act.Tags().MatchAllTag(tag.ContainerFromTag(tags.Attack, tags.Melee, tags.Weapon)) {
			continue
		}
		if act.AverageDamage() > bestDamage {
			best = act
		}
	}
	return best
}
