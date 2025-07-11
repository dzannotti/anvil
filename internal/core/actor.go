package core

import (
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Actor struct {
	Dispatcher         EventDispatcher
	Encounter          *Encounter
	Position           grid.Position
	World              *World
	Attributes         stats.Attributes
	Proficiencies      stats.Proficiencies
	SpellCastingSource tag.Tag
	Name               string
	HitPoints          int
	MaxHitPoints       int
	Actions            []Action
	Team               TeamID
	Effects            EffectContainer
	Equipped           []Item
	Resources          Resources
	Conditions         Conditions
}

func (a *Actor) StartTurn() {
	a.Resources.Reset()
	a.Evaluate(&TurnStarted{Source: a})
}

func (a *Actor) EndTurn() {
	a.Evaluate(&TurnEnded{Source: a})
}

func (a *Actor) Evaluate(state any) {
	a.Effects.Evaluate(state)
}

func (a *Actor) AddAction(action ...Action) {
	for _, ca := range action {
		if a.HasAction(ca) {
			continue
		}

		a.Actions = append(a.Actions, ca)
	}
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

func (a *Actor) AddCondition(t tag.Tag, src *Effect) {
	a.Conditions.Add(t, src)
	a.Evaluate(&ConditionChanged{Source: a, From: src, Condition: t})
	a.Dispatcher.Emit(ConditionChangedEvent{Source: a, From: src, Condition: t, Added: true})
}

func (a *Actor) RemoveCondition(t tag.Tag, src *Effect) {
	ok := a.Conditions.Remove(t, src)
	if !ok {
		return
	}

	a.Evaluate(&ConditionChanged{Source: a, From: src, Condition: t})
	a.Dispatcher.Emit(ConditionChangedEvent{Source: a, From: src, Condition: t, Added: false})
}

func (a *Actor) Equip(item Item) {
	a.Equipped = append(a.Equipped, item)
	item.OnEquip(a)
}

func (a *Actor) Die() {
	a.Dispatcher.Begin(DeathEvent{Actor: a})
	defer a.Dispatcher.End()
	a.AddCondition(tags.Dead, &Effect{Name: "Dead"})
	a.Dispatcher.Emit(ConfirmEvent{Confirm: true})
}

func (a *Actor) ConsumeResource(t tag.Tag, amount int) {
	a.Resources.Consume(t, amount)
	a.Dispatcher.Emit(SpendResourceEvent{Source: a, Resource: t, Amount: amount})
}
