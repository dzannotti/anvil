package team

import "anvil/internal/core/creature"

func (t *Team) AddMember(creature *creature.Creature) {
	t.members = append(t.members, creature)
}

func (t *Team) Members() []*creature.Creature {
	return t.members
}
