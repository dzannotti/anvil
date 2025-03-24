package encounter

import (
	"sync"

	"anvil/internal/core/definition"
	"anvil/internal/core/event"
)

type Act = func(active definition.Creature, wg *sync.WaitGroup)

func (e *Encounter) playTurn(act Act) {
	//log.Start(NewTurnEvent(turn+1, active))
	//defer log.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	e.ActiveCreature().StartTurn()
	go act(e.ActiveCreature(), &turnWG)
	turnWG.Wait()
}

func (e *Encounter) playRound(act Act) {
	//log.Start(NewRoundEvent(round+1, creatures))
	//defer log.End()
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
	e.log.Start(event.NewEncounter(e.teams))
	defer e.log.End()
	defer wg.Done()
	for !e.IsOver() {
		e.playRound(act)
		e.round = e.round + 1
		//log.End() // round
	}
}
