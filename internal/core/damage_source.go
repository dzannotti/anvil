package core

import "anvil/internal/tag"

type DamageSource struct {
	Times  int
	Sides  int
	Source string
	Tags   tag.Container
}
