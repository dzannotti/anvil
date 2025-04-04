package weapon

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func NewDagger() Weapon {
	return Weapon{
		name:   "Dagger",
		damage: []core.DamageSource{{Times: 1, Sides: 4, Source: "Dagger", Tags: tag.ContainerFromTags([]tag.Tag{tags.Piercing})}},
	}
}

func NewGreatAxe() Weapon {
	return Weapon{
		name:   "Great Axe",
		damage: []core.DamageSource{{Times: 2, Sides: 6, Source: "Great Axe", Tags: tag.ContainerFromTags([]tag.Tag{tags.Slashing})}},
	}
}
