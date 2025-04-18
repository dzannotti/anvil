package core

func (e Encounter) IsOver() bool {
	alive := 0
	teams := []TeamID{TeamPlayers, TeamEnemies}
	for _, t := range teams {
		if !e.IsTeamDead(t) {
			alive = alive + 1
		}
	}
	return alive <= 1
}

func (e Encounter) IsTeamDead(team TeamID) bool {
	for _, c := range e.Actors {
		if c.Team == team && !c.IsDead() {
			return false
		}
	}
	return true
}

func (e Encounter) ActiveActor() *Actor {
	return e.InitiativeOrder[e.Turn]
}

func (e Encounter) Winner() (TeamID, bool) {
	for _, c := range e.Actors {
		if !c.IsDead() {
			return c.Team, true
		}
	}
	return "", false
}
