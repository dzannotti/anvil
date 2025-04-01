package core

import (
	"anvil/internal/core/stats"
	"anvil/internal/effect"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
)

type Creature struct {
	Log           *eventbus.Hub
	Position      grid.Position
	World         *World
	Attributes    stats.Attributes
	Proficiencies stats.Proficiencies
	Name          string
	HitPoints     int
	MaxHitPoints  int
	Actions       []Action
	Effects       effect.Container
}

func (c *Creature) Evaluate(event string, state any) {
	c.Effects.Evaluate(event, state)
}

func (c Creature) IsDead() bool {
	return c.HitPoints == 0
}

func (c *Creature) AddAction(action ...Action) {
	c.Actions = append(c.Actions, action...)
}
