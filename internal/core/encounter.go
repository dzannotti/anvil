package core

import (
	"slices"
	"sync"
)

type Encounter struct {
	Round           int
	Turn            int
	InitiativeOrder []*Actor
	Actors          []*Actor
	Log             LogWriter
	World           *World
}

type Act = func(active *Actor)

func (e *Encounter) playTurn(act Act) {
	e.Log.Start(TurnType, TurnEvent{Turn: e.Turn, Actor: e.ActiveActor()})
	defer e.Log.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	e.ActiveActor().StartTurn()
	go func() {
		defer turnWG.Done()
		act(e.ActiveActor())
	}()
	e.ActiveActor().EndTurn()
	turnWG.Wait()
}

func (e *Encounter) EndTurn(act Act) {
	e.ActiveActor().EndTurn()

	e.Log.Start(TurnType, TurnEvent{Turn: e.Turn, Actor: e.ActiveActor()})
	defer e.Log.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	e.ActiveActor().StartTurn()
	go func() {
		defer turnWG.Done()
		act(e.ActiveActor())
	}()

	turnWG.Wait()
}

func (e *Encounter) playRound(act Act) {
	e.Log.Start(RoundType, RoundEvent{Round: e.Round, Actors: e.Actors})
	defer e.Log.End()
	e.Turn = 0
	for i := range e.InitiativeOrder {
		e.Turn = i
		e.playTurn(act)
		if e.IsOver() {
			break
		}
	}
}

func (e *Encounter) Start() {
	e.Round = 0
	e.Turn = 0
	for _, a := range e.Actors {
		a.Encounter = e
	}
	e.InitiativeOrder = slices.Clone(e.Actors)
	e.Log.Start(EncounterType, EncounterEvent{Actors: e.Actors, World: e.World})
}

func (e *Encounter) Play(act Act) string {
	e.Round = 0
	e.Turn = 0
	for _, a := range e.Actors {
		a.Encounter = e
	}
	e.InitiativeOrder = slices.Clone(e.Actors)
	e.Log.Start(EncounterType, EncounterEvent{Actors: e.Actors, World: e.World})
	defer e.Log.End()
	for !e.IsOver() {
		e.playRound(act)
		e.Round = e.Round + 1
	}
	winner, _ := e.Winner()
	return winner
}
