package core

import (
	"anvil/internal/core/definition"
)

func (t *Team) AddMember(creature definition.Creature) {
	t.members = append(t.members, creature)
}

func (t *Team) Members() []definition.Creature {
	return t.members
}
