package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type AttackRoll struct {
	Source snapshot.Creature
	Target snapshot.Creature
}

func NewAttackRoll(source definition.Creature, target definition.Creature) (string, AttackRoll) {
	return "attack_roll", AttackRoll{
		Source: snapshot.CaptureCreature(source),
		Target: snapshot.CaptureCreature(target),
	}
}
