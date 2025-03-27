package encounter

import (
	"anvil/internal/core/definition"
	"anvil/internal/log"
)

type Encounter struct {
	round           int
	turn            int
	initiativeOrder []definition.Creature
	teams           []definition.Team
	log             *log.EventLog
	world           definition.World
}

func New(log *log.EventLog, world definition.World, teams []definition.Team) *Encounter {
	encounter := &Encounter{
		world:           world,
		log:             log,
		teams:           teams,
		initiativeOrder: []definition.Creature{},
	}
	encounter.initiativeOrder = encounter.AllCreatures()
	return encounter
}
