package team

import "anvil/internal/core/creature"

type Team struct {
	name    string
	members []*creature.Creature
}

func New(name string) *Team {
	return &Team{
		name:    name,
		members: []*creature.Creature{},
	}
}
