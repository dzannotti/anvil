package core

import (
	"anvil/internal/core/stats"
	"anvil/internal/effect"
	"anvil/internal/grid"
)

type Actor struct {
	Log           LogWriter
	Position      grid.Position
	World         *World
	Attributes    stats.Attributes
	Proficiencies stats.Proficiencies
	Name          string
	HitPoints     int
	MaxHitPoints  int
	Actions       []Action
	Team          TeamID
	Effects       effect.Container
}

func (a *Actor) StartTurn() {

}

func (a *Actor) Evaluate(event string, state any) {
	a.Effects.Evaluate(event, state)
}

func (a *Actor) AddAction(action ...Action) {
	a.Actions = append(a.Actions, action...)
}

func (a Actor) IsDead() bool {
	return a.HitPoints == 0
}
