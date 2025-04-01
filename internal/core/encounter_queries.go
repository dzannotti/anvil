package core

func (e Encounter) IsOver() bool {
	alive := 0
	teams := []TeamID{TeamPlayers, TeamEnemies, TeamNeutral, TeamGaea}
	for _, t := range teams {
		if !e.IsTeamDead(t) {
			alive = alive + 1
		}
	}
	return alive <= 1
}

func (e Encounter) IsTeamDead(team TeamID) bool {
	for _, c := range e.Creatures {
		if c.Team == team && !c.IsDead() {
			return false
		}
	}
	return true
}

func (e Encounter) ActiveCreature() *Creature {
	return e.InitiativeOrder[e.Turn]
}

func (e Encounter) Winner() (string, bool) {
	for _, c := range e.Creatures {
		if !c.IsDead() {
			return c.Name, true
		}
	}
	return "", false
}
