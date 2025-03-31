package core

import "anvil/internal/core/definition"

func (t Team) IsDead() bool {
	for _, c := range t.members {
		if !c.IsDead() {
			return false
		}
	}
	return true
}

func (t Team) Name() string {
	return t.name
}

func (t Team) Contains(c definition.Creature) bool {
	for _, m := range t.members {
		if m == c {
			return true
		}
	}
	return false
}
