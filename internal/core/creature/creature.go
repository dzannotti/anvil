package creature

import (
	"anvil/internal/core/definition"
	"anvil/internal/grid"
	"anvil/internal/log"
)

type Creature struct {
	log          *log.EventLog
	position     grid.Position
	world        definition.World
	name         string
	hitPoints    int
	maxHitPoints int
	actions      []definition.Action
}

func New(log *log.EventLog, world definition.World, pos grid.Position, name string, hitPoints int) *Creature {
	return &Creature{
		log:          log,
		position:     pos,
		world:        world,
		name:         name,
		hitPoints:    hitPoints,
		maxHitPoints: hitPoints,
	}
}
