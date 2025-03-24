package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/parts"
)

type TakeDamage struct {
	Target       parts.Creature
	Damage       int
	HitPoints    int
	MaxHitPoints int
}

func (e TakeDamage) Type() string {
	return "CreatureTakeDamage"
}

func NewTakeDamage(target definition.Creature, damage int) TakeDamage {
	return TakeDamage{Target: parts.NewCreature(target), Damage: damage, HitPoints: target.HitPoints(), MaxHitPoints: target.MaxHitPoints()}
}
