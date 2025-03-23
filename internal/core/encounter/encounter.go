package encounter

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/team"
	"anvil/internal/log"
	"sync"
)

type Act = func(active *creature.Creature, creatures []*creature.Creature, wg *sync.WaitGroup)

func IsOver(creatures []*creature.Creature) bool {
	playersAlive := false
	enemiesAlive := false
	for _, c := range creatures {
		if !c.IsDead() {
			if c.FactionID() == team.Player {
				playersAlive = true
			}
			if c.FactionID() == team.Enemy {
				enemiesAlive = true
			}
		}
	}
	return !playersAlive || !enemiesAlive
}

func winner(creatures []*creature.Creature) team.Team {
	for i := range creatures {
		if !creatures[i].IsDead() {
			return creatures[i].FactionID()
		}
	}
	return team.None
}

func Play(log *log.EventLog, creatures []*creature.Creature, act Act, result chan team.Team) {
	turn := 0
	round := 0
	log.Start(NewEncounterEvent(creatures))
	defer log.End()
	for !IsOver(creatures) {
		log.Start(NewRoundEvent(round+1, creatures))
		for i := range creatures {
			var active = creatures[i]
			log.Start(NewTurnEvent(turn+1, active))
			turnWG := sync.WaitGroup{}
			turnWG.Add(1)
			active.StartTurn()
			go act(active, creatures, &turnWG)
			turnWG.Wait()
			log.End() // turn
			turn = turn + 1
			if IsOver(creatures) {
				break
			}
		}
		round = round + 1
		turn = 0
		log.End() // round
	}
	result <- winner(creatures)
}
