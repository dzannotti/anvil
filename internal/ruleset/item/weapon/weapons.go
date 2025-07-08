package weapon

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

func NewDagger() *Weapon {
	return &Weapon{
		archetype: "dagger",
		id:        uuid.New().String(),
		name:      "Dagger",
		tags:      tag.NewContainer(tags.Melee, tags.SimpleWeapon),
		damage:    []core.DamageSource{core.NewLegacyDamageSource(1, 4, "Dagger", tag.NewContainer(tags.Piercing))},
		reach:     1,
	}
}

func NewGreatAxe() *Weapon {
	return &Weapon{
		archetype: "great-axe",
		id:        uuid.New().String(),
		name:      "Great Axe",
		tags:      tag.NewContainer(tags.Melee, tags.MartialAxe),
		damage: []core.DamageSource{
			core.NewLegacyDamageSource(2, 6, "Great Axe", tag.NewContainer(tags.Slashing)),
		},
		reach: 1,
	}
}
