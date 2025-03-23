package encounter

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/team"
)

func IsOver(creatures []*creature.Creature) bool {
	playersAlive := false
	enemiesAlive := false
	for _, c := range creatures {
		if !c.IsDead() {
			if c.Team() == team.Player {
				playersAlive = true
			}
			if c.Team() == team.Enemy {
				enemiesAlive = true
			}
		}
	}
	return !playersAlive || !enemiesAlive
}

func winner(creatures []*creature.Creature) team.Team {
	for i := range creatures {
		if !creatures[i].IsDead() {
			return creatures[i].Team()
		}
	}
	return team.None
}
