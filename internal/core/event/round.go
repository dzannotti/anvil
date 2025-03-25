package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type Round struct {
	Round     int
	Creatures []snapshot.Creature
}

func NewRound(round int, c []definition.Creature) Round {
	creatures := make([]snapshot.Creature, 0, len(c))
	for i := range c {
		creatures = append(creatures, snapshot.CaptureCreature(c[i]))
	}
	return Round{Round: round, Creatures: creatures}
}
