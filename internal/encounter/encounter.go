package encounter

import (
	"anvil/internal/creature"
	"anvil/internal/team"
	"fmt"
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

func Play(allCreatures []*creature.Creature, act Act, result chan team.Team) {
	turn := 0
	round := 0
	for !IsOver(allCreatures) {
		fmt.Println("Round", round+1)
		for i := range allCreatures {
			var activeCreature = allCreatures[i]
			fmt.Println("Turn", turn+1, activeCreature.Name(), "turn")
			turnWG := sync.WaitGroup{}
			turnWG.Add(1)
			activeCreature.StartTurn()
			go act(activeCreature, allCreatures, &turnWG)
			turnWG.Wait()
			turn = turn + 1
			if IsOver(allCreatures) {
				break
			}
		}
		round = round + 1
		turn = 0
	}
	result <- winner(allCreatures)
}
