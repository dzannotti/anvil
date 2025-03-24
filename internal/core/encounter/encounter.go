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
	log             log.EventLog
}

func New(log log.EventLog, teams []definition.Team) *Encounter {
	encounter := &Encounter{
		teams:           teams,
		initiativeOrder: []definition.Creature{},
		log:             log,
	}
	encounter.initiativeOrder = encounter.AllCreatures()
	return encounter
}
