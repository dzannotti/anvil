package encounter

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/team"
	"anvil/internal/log"
	"sync"
)

type Act = func(active *creature.Creature, creatures []*creature.Creature, wg *sync.WaitGroup)

func playTurn(log *log.EventLog, turn int, active *creature.Creature, creatures []*creature.Creature, act Act) {
	log.Start(NewTurnEvent(turn+1, active))
	defer log.End()
	turnWG := sync.WaitGroup{}
	turnWG.Add(1)
	active.StartTurn()
	go act(active, creatures, &turnWG)
	turnWG.Wait()
}

func playRound(log *log.EventLog, round int, creatures []*creature.Creature, act Act) {
	log.Start(NewRoundEvent(round+1, creatures))
	defer log.End()
	turn := 0
	for i := range creatures {
		var active = creatures[i]
		turn = turn + 1
		playTurn(log, turn, active, creatures, act)
		if IsOver(creatures) {
			break
		}
	}
}

func Play(log *log.EventLog, creatures []*creature.Creature, act Act, result chan team.Team) {
	round := 0
	log.Start(NewEncounterEvent(creatures))
	defer log.End()
	for !IsOver(creatures) {
		playRound(log, round, creatures, act)
		round = round + 1
		log.End() // round
	}
	result <- winner(creatures)
}
