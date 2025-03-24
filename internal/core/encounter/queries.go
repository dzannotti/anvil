package encounter

import (
	"anvil/internal/core/definition"
)

func (e Encounter) IsOver() bool {
	alive := 0
	for _, t := range e.teams {
		if !t.IsDead() {
			alive++
		}
	}
	return alive <= 1
}

func (e Encounter) ActiveCreature() definition.Creature {
	return e.initiativeOrder[e.turn]
}

func (e Encounter) AllCreatures() []definition.Creature {
	var allCreatures = []definition.Creature{}
	for _, t := range e.teams {
		allCreatures = append(allCreatures, t.Members()...)
	}
	return allCreatures
}

func (e Encounter) Winner() definition.Team {
	for _, t := range e.teams {
		if !t.IsDead() {
			return t
		}
	}
	return nil
}

func (e Encounter) Teams() []definition.Team {
	return e.teams
}
