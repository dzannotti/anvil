package core

import "anvil/internal/core/definition"

func (t Team) IsDead() bool {
	for _, c := range t.Members {
		if !c.IsDead() {
			return false
		}
	}
	return true
}

func (t Team) Contains(c definition.Creature) bool {
	for _, m := range t.Members {
		if m == c {
			return true
		}
	}
	return false
}
