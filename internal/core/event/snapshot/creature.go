package snapshot

import "anvil/internal/core/definition"

type Creature struct {
	Name         string
	HitPoints    int
	MaxHitPoints int
}

func CaptureCreature(src definition.Creature) Creature {
	return Creature{Name: src.Name(), HitPoints: src.HitPoints(), MaxHitPoints: src.MaxHitPoints()}
}