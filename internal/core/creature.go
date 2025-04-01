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
	Effects       effect.Container
}

func (c *Creature) Evaluate(event string, state any) {
	c.Effects.Evaluate(event, state)
}
