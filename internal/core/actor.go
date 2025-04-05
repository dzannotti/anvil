package core

import (
	"anvil/internal/core/stats"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Actor struct {
	Log           LogWriter
	Encounter     *Encounter
	Position      grid.Position
	World         *World
	Attributes    stats.Attributes
	Proficiencies stats.Proficiencies
	Name          string
	HitPoints     int
	MaxHitPoints  int
	Actions       []Action
	Team          TeamID
	Effects       EffectContainer
	Equipped      []Item
	Resources     Resources
}

func (a *Actor) StartTurn() {

}

func (a *Actor) Evaluate(event string, state any) {
	a.Effects.Evaluate(event, state)
}

func (a *Actor) AddAction(action ...Action) {
	a.Actions = append(a.Actions, action...)
}

func (a *Actor) AddEffect(effect ...*Effect) {
	a.Effects.Add(effect...)
}

func (a *Actor) RemoveEffect(effect *Effect) {
	a.Effects.Remove(effect)
}

func (a *Actor) AddProficiency(t tag.Tag) {
	a.Proficiencies.Add(t)
}

func (a Actor) IsDead() bool {
	return a.HitPoints == 0
}

func (a Actor) CanAct() bool {
	return !a.IsDead()
}

func (a *Actor) Equip(item Item) {
	a.Equipped = append(a.Equipped, item)
	item.OnEquip(a)
}
