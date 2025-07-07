package weapon

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func NewDagger() *Weapon {
	return &Weapon{
		name:   "Dagger",
		tags:   tag.NewContainer(tags.Melee, tags.SimpleWeapon),
		damage: []core.DamageSource{{Times: 1, Sides: 4, Source: "Dagger", Tags: tag.NewContainer(tags.Piercing)}},
		reach:  1,
	}
}

func NewGreatAxe() *Weapon {
	return &Weapon{
		name: "Great Axe",
		tags: tag.NewContainer(tags.Melee, tags.MartialAxe),
		damage: []core.DamageSource{
			{Times: 2, Sides: 6, Source: "Great Axe", Tags: tag.NewContainer(tags.Slashing)},
		},
		reach: 1,
	}
}
