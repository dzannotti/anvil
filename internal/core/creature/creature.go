package creature

import (
	"anvil/internal/core/definition"
	"anvil/internal/effect"
	"anvil/internal/effect/state"
	"anvil/internal/grid"
	"anvil/internal/log"
)

type Creature struct {
	log           *log.EventLog
	position      grid.Position
	world         definition.World
	attributes    Attributes
	proficiencies Proficiencies
	name          string
	hitPoints     int
	maxHitPoints  int
	actions       []definition.Action
	effects       *effect.Container
}

func New(log *log.EventLog, world definition.World, pos grid.Position, name string, hitPoints int, attributes Attributes, proficiencies Proficiencies) *Creature {
	creature := &Creature{
		log:           log,
		position:      pos,
		world:         world,
		name:          name,
		effects:       effect.NewContainer(),
		hitPoints:     hitPoints,
		maxHitPoints:  hitPoints,
		attributes:    attributes,
		proficiencies: proficiencies,
	}
	cell, _ := world.At(pos)
	cell.AddOccupant(creature)
	return creature
}

func (c *Creature) Evaluate(state state.State) {
	c.effects.Evaluate(state)
}
