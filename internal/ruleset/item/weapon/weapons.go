package weapon

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func NewDagger() Weapon {
	return Weapon{
		name:   "Dagger",
		damage: []core.DamageSource{{Times: 1, Sides: 4, Source: "Dagger", Tags: tag.ContainerFromTag(tags.Piercing)}},
		reach:  10,
	}
}

func NewGreatAxe() Weapon {
	return Weapon{
		name:   "Great Axe",
		damage: []core.DamageSource{{Times: 2, Sides: 6, Source: "Great Axe", Tags: tag.ContainerFromTag(tags.Slashing)}},
		reach:  10,
	}
}
