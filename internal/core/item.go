package core

import "anvil/internal/tag"

type Item interface {
	OnEquip(a *Actor)
	Tags() tag.Container
}
