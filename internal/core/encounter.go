package core

import (
	"slices"
)

type Encounter struct {
	Round           int
	Turn            int
	InitiativeOrder []*Actor
	Actors          []*Actor
	Log             LogWriter
	World           *World
}

func (e *Encounter) Start() {
	for _, a := range e.Actors {
		a.Encounter = e
	}
	e.InitiativeOrder = slices.Clone(e.Actors)
	e.Log.Start(EncounterType, EncounterEvent{Actors: e.Actors, World: e.World})
	e.Round = -1
	e.startRound()
	e.startTurn()
}

func (e *Encounter) End() {
	if !e.IsOver() {
		return
	}
	e.Log.End()
}

func (e *Encounter) EndTurn() {
	e.ActiveActor().EndTurn()
	e.Log.End()
	if e.IsOver() {
		e.endRound()
		return
	}
	e.Turn = e.Turn + 1
	if e.Turn >= len(e.InitiativeOrder) {
		e.endRound()
		e.startRound()
	}
	e.startTurn()
}

func (e *Encounter) startRound() {
	e.Round = e.Round + 1
	e.Log.Start(RoundType, RoundEvent{Round: e.Round, Actors: e.Actors})
	e.Turn = 0
}

func (e *Encounter) endRound() {
	e.Log.End()
}

func (e *Encounter) startTurn() {
	e.Log.Start(TurnType, TurnEvent{Turn: e.Turn, Actor: e.ActiveActor()})
	e.ActiveActor().StartTurn()
}
