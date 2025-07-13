package core

import (
	"slices"
)

type Encounter struct {
	Round           int
	Turn            int
	InitiativeOrder []*Actor
	Actors          []*Actor
	Dispatcher      EventDispatcher
	World           *World
}

func (e *Encounter) Start() {
	for _, a := range e.Actors {
		a.Encounter = e
	}
	e.InitiativeOrder = slices.Clone(e.Actors)
	e.Dispatcher.Begin(EncounterEvent{Actors: e.Actors, World: e.World})
	e.Round = -1
	e.startRound()
	e.startTurn()
}

func (e *Encounter) End() {
	// This method is now mostly a no-op since EndTurn() handles
	// ending the encounter when it's over. We keep it for compatibility
	// but make it safe to call multiple times.
}

func (e *Encounter) EndTurn() {
	e.ActiveActor().EndTurn()
	e.Dispatcher.End()
	if e.IsOver() {
		e.endRound()
		e.Dispatcher.End() // End the encounter when it's over
		return
	}
	e.Turn++
	if e.Turn >= len(e.InitiativeOrder) {
		e.endRound()
		e.startRound()
	}
	e.startTurn()
}

func (e *Encounter) startRound() {
	e.Round++
	e.Dispatcher.Begin(RoundEvent{Round: e.Round, Actors: e.Actors})
	e.Turn = 0
}

func (e *Encounter) endRound() {
	e.Dispatcher.End()
}

func (e *Encounter) startTurn() {
	e.Dispatcher.Begin(TurnEvent{Turn: e.Turn, Actor: e.ActiveActor()})
	e.ActiveActor().StartTurn()
}
