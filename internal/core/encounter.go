package core

import (
	"anvil/internal/eventbus"
	"sync"
)

type Encounter struct {
	Round           int
	Turn            int
	InitiativeOrder []*Creature
	Teams           []*Team
	Hub             *eventbus.Hub
	World           *World
}

type Act = func(active *Creature, wg *sync.WaitGroup)

func (e *Encounter) playTurn(act Act) {
	e.Hub.Start(TurnEventType, TurnEvent{Turn: e.Turn, Creature: *e.ActiveCreature()})
	defer e.Hub.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	e.ActiveCreature().StartTurn()
	go act(e.ActiveCreature(), &turnWG)
	turnWG.Wait()
}

func (e *Encounter) playRound(act Act) {
	e.Hub.Start(RoundEventType, RoundEvent{Round: e.Round, Creatures: e.AllCreatures()})
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
	e.InitiativeOrder = e.AllCreatures()
	e.Hub.Start(EncounterEventType, EncounterEvent{Teams: e.Teams, World: *e.World})
	defer e.Hub.End()
	defer wg.Done()
	for !e.IsOver() {
		e.playRound(act)
		e.Round = e.Round + 1
	}
}
