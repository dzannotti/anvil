package core

func (e Encounter) IsOver() bool {
	alive := 0
	for _, t := range e.teams {
		if !t.IsDead() {
			alive++
		}
	}
	return alive <= 1
}

func (e Encounter) ActiveCreature() *Creature {
	return e.initiativeOrder[e.turn]
}

func (e Encounter) AllCreatures() []*Creature {
	var allCreatures = []*Creature{}
	for _, t := range e.teams {
		allCreatures = append(allCreatures, t.Members...)
	}
	return allCreatures
}

func (e Encounter) Winner() (Team, bool) {
	for _, t := range e.teams {
		if !t.IsDead() {
			return *t, true
		}
	}
	return Team{}, false
}

func (e Encounter) Teams() []*Team {
	return e.teams
}
