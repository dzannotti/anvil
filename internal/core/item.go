package core

import "anvil/internal/tag"

type Item interface {
	Name() string
	OnEquip(a *Actor)
	Tags() tag.Container
}
