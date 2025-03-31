package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/eventbus"
)

type Encounter struct {
	round           int
	turn            int
	initiativeOrder []definition.Creature
	teams           []definition.Team
	hub             *eventbus.Hub
	world           definition.World
}

func NewEncounter(hub *eventbus.Hub, world definition.World, teams []definition.Team) *Encounter {
	encounter := &Encounter{
		world:           world,
		hub:             hub,
		teams:           teams,
		initiativeOrder: []definition.Creature{},
	}
	encounter.initiativeOrder = encounter.AllCreatures()
	return encounter
}
