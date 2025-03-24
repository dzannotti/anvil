package parts

import "anvil/internal/core/definition"

type Creature struct {
	Name         string
	HitPoints    int
	MaxHitPoints int
}

func NewCreature(src definition.Creature) Creature {
	return Creature{Name: src.Name(), HitPoints: src.HitPoints(), MaxHitPoints: src.MaxHitPoints()}
}
