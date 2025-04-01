package core

func (e Encounter) IsOver() bool {
	alive := 0
	for _, t := range e.Teams {
		if !t.IsDead() {
			alive++
		}
	}
	return alive <= 1
}

func (e Encounter) ActiveCreature() *Creature {
	return e.InitiativeOrder[e.Turn]
}

func (e Encounter) AllCreatures() []*Creature {
	var allCreatures = []*Creature{}
	for _, t := range e.Teams {
		allCreatures = append(allCreatures, t.Members...)
	}
	return allCreatures
}

func (e Encounter) Winner() (Team, bool) {
	for _, t := range e.Teams {
		if !t.IsDead() {
			return *t, true
		}
	}
	return Team{}, false
}
