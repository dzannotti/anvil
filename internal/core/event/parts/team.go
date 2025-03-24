package parts

import "anvil/internal/core/definition"

type Team struct {
	Name    string
	Members []Creature
}

func NewFaction(team definition.Team) Team {
	members := make([]Creature, len(team.Members()))
	for i, member := range team.Members() {
		members[i] = NewCreature(member)
	}
	return Team{Name: team.Name(), Members: members}
}
