package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/effect"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
)

type Creature struct {
	log           *eventbus.Hub
	position      grid.Position
	world         *World
	attributes    Attributes
	proficiencies Proficiencies
	name          string
	hitPoints     int
	maxHitPoints  int
	actions       []definition.Action
	effects       *effect.Container
}

func NewCreature(log *eventbus.Hub, world *World, pos grid.Position, name string, hitPoints int, attributes Attributes, proficiencies Proficiencies) *Creature {
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
	return creature
}

func (c *Creature) Evaluate(event string, state any) {
	c.effects.Evaluate(event, state)
}
