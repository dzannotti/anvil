package core

import (
	"anvil/internal/eventbus"
	"slices"
	"sync"
)

type Encounter struct {
	Round           int
	Turn            int
	InitiativeOrder []*Actor
	Actors          []*Actor
	Hub             *eventbus.Hub
	World           *World
}

type Act = func(active *Actor, wg *sync.WaitGroup)

func (e *Encounter) playTurn(act Act) {
	e.Hub.Start(TurnEventType, TurnEvent{Turn: e.Turn, Actor: *e.ActiveActor()})
	defer e.Hub.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	e.ActiveActor().StartTurn()
	go act(e.ActiveActor(), &turnWG)
	turnWG.Wait()
}

func (e *Encounter) playRound(act Act) {
	e.Hub.Start(RoundEventType, RoundEvent{Round: e.Round, Actors: e.Actors})
	defer e.Hub.End()
	e.Turn = 0
	for i := range e.InitiativeOrder {
		e.Turn = i
		e.playTurn(act)
		if e.IsOver() {
			break
		}
	}
}

func (e *Encounter) Play(act Act, wg *sync.WaitGroup) {
	e.Round = 0
	e.Turn = 0
	e.InitiativeOrder = slices.Clone(e.Actors)
	e.Hub.Start(EncounterEventType, EncounterEvent{Actors: e.Actors, World: *e.World})
	defer e.Hub.End()
	defer wg.Done()
	for !e.IsOver() {
		e.playRound(act)
		e.Round = e.Round + 1
	}
}
