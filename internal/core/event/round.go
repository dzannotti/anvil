package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/parts"
)

type Round struct {
	Round     int
	Creatures []parts.Creature
}

func NewRound(round int, c []definition.Creature) Round {
	creatures := make([]parts.Creature, 0, len(c))
	for i := range c {
		creatures = append(creatures, parts.NewCreature(c[i]))
	}
	return Round{Round: round, Creatures: creatures}
}
