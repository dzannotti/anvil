package core

import "anvil/internal/tag"

type Item interface {
	Name() string
	Archetype() string
	ID() string
	OnEquip(a *Actor)
	Tags() *tag.Container
}
