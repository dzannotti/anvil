package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type Died struct {
	Creature snapshot.Creature
}

func NewDied(src definition.Creature) Died {
	return Died{Creature: snapshot.CaptureCreature(src)}
}
