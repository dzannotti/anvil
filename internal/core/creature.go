package core

import (
	"anvil/internal/effect"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
)

type Creature struct {
	Log           *eventbus.Hub
	Position      grid.Position
	World         *World
	Attributes    Attributes
	Proficiencies Proficiencies
	Name          string
	HitPoints     int
	MaxHitPoints  int
	Actions       []Action
	Effects       *effect.Container
}

func NewCreature(log *eventbus.Hub, world *World, pos grid.Position, name string, hitPoints int, attributes Attributes, proficiencies Proficiencies) *Creature {
	creature := &Creature{
		Log:           log,
		Position:      pos,
		World:         world,
		Name:          name,
		Effects:       effect.NewContainer(),
		HitPoints:     hitPoints,
		MaxHitPoints:  hitPoints,
		Attributes:    attributes,
		Proficiencies: proficiencies,
	}
	return creature
}

func (c *Creature) Evaluate(event string, state any) {
	c.Effects.Evaluate(event, state)
}
