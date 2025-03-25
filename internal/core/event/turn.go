package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type Turn struct {
	Turn     int
	Creature snapshot.Creature
}

func NewTurn(turn int, src definition.Creature) Turn {
	return Turn{Turn: turn, Creature: snapshot.CaptureCreature(src)}
}
