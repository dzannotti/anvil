package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/parts"
)

type Turn struct {
	Turn     int
	Creature parts.Creature
}

func NewTurn(turn int, src definition.Creature) Turn {
	return Turn{Turn: turn, Creature: parts.NewCreature(src)}
}
