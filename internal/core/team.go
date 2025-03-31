package core

import (
	"anvil/internal/core/definition"
)

type Team struct {
	name    string
	members []definition.Creature
}

func NewTeam(name string) *Team {
	return &Team{
		name:    name,
		members: []definition.Creature{},
	}
}
