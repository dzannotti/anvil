package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/parts"
)

type Died struct {
	Creature parts.Creature
}

func NewDied(src definition.Creature) Died {
	return Died{Creature: parts.NewCreature(src)}
}
