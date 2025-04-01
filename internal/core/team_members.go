package core

import (
	"anvil/internal/core/definition"
)

func (t *Team) AddMember(creature definition.Creature) {
	t.Members = append(t.Members, creature)
}
