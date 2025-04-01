package core

import (
	"sync"

	"anvil/internal/core/definition"
)

type Act = func(active definition.Creature, wg *sync.WaitGroup)

func (e *Encounter) playTurn(act Act) {
	e.hub.Start(NewTurnEvent(e.turn, e.ActiveCreature()))
	defer e.hub.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	e.ActiveCreature().StartTurn()
	go act(e.ActiveCreature(), &turnWG)
	turnWG.Wait()
}

func (e *Encounter) playRound(act Act) {
	e.hub.Start(NewRoundEvent(e.round, e.AllCreatures()))
	defer e.hub.End()
	e.turn = 0
	for i := range e.initiativeOrder {
		e.turn = i
		e.playTurn(act)
		if e.IsOver() {
			break
		}
	}
}

func (e *Encounter) Play(act Act, wg *sync.WaitGroup) {
	e.round = 0
	e.hub.Start(NewEncounterEvent(e))
	defer e.hub.End()
	defer wg.Done()
	for !e.IsOver() {
		e.playRound(act)
		e.round = e.round + 1
	}
}
