package core

import (
	"anvil/internal/eventbus"
)

type Encounter struct {
	round           int
	turn            int
	initiativeOrder []*Creature
	teams           []*Team
	hub             *eventbus.Hub
	World           *World
}

func NewEncounter(hub *eventbus.Hub, world *World, teams []*Team) *Encounter {
	encounter := &Encounter{
		World:           world,
		hub:             hub,
		teams:           teams,
		initiativeOrder: []*Creature{},
	}
	encounter.initiativeOrder = encounter.AllCreatures()
	return encounter
}
