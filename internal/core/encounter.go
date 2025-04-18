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
	e.Round = 0
	e.Turn = 0
	for _, a := range e.Actors {
		a.Encounter = e
	}
	e.InitiativeOrder = slices.Clone(e.Actors)
	e.Log.Start(EncounterType, EncounterEvent{Actors: e.Actors, World: e.World})
	e.Log.Start(RoundType, RoundEvent{Round: e.Round, Actors: e.Actors})
	e.Log.Start(TurnType, TurnEvent{Turn: e.Turn, Actor: e.ActiveActor()})
}

func (e *Encounter) End() TeamID {
	if !e.IsOver() {
		panic("encounter is not over")
	}
	e.Log.End()
	winner, _ := e.Winner()
	return winner
}

func (e *Encounter) EndTurn() {
	e.Log.End()
	e.Turn = e.Turn + 1
	if e.Turn >= len(e.InitiativeOrder) {
		e.EndRound()
	}
	if e.IsOver() {
		e.EndRound()
		return
	}
	e.Log.Start(TurnType, TurnEvent{Turn: e.Turn, Actor: e.ActiveActor()})
	e.ActiveActor().StartTurn()
}

func (e *Encounter) EndRound() {
	e.Log.End()
	if e.IsOver() {
		return
	}
	e.Round = e.Round + 1
	e.Log.Start(RoundType, RoundEvent{Round: e.Round, Actors: e.Actors})
	e.Turn = 0
}
