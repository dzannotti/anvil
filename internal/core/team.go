package core

import (
	"anvil/internal/core/definition"
)

type Team struct {
	Name    string
	Members []definition.Creature
}

func NewTeam(name string) *Team {
	return &Team{
		Name:    name,
		Members: []definition.Creature{},
	}
}
