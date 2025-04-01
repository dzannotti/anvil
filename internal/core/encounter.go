package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/eventbus"
)

type Encounter struct {
	round           int
	turn            int
	initiativeOrder []definition.Creature
	teams           []*Team
	hub             *eventbus.Hub
	World           *World
}

func NewEncounter(hub *eventbus.Hub, world *World, teams []*Team) *Encounter {
	encounter := &Encounter{
		World:           world,
		hub:             hub,
		teams:           teams,
		initiativeOrder: []definition.Creature{},
	}
	encounter.initiativeOrder = encounter.AllCreatures()
	return encounter
}
