package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type TakeDamage struct {
	Target       snapshot.Creature
	Damage       int
	HitPoints    int
	MaxHitPoints int
}

func NewTakeDamage(target definition.Creature, damage int) (string, TakeDamage) {
	return "take_damage", TakeDamage{Target: snapshot.CaptureCreature(target), Damage: damage, HitPoints: target.HitPoints(), MaxHitPoints: target.MaxHitPoints()}
}
