package weapon

import (
	"anvil/internal/core/tags"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

func NewDagger() *Weapon {
	return &Weapon{
		archetype:   "dagger",
		id:          uuid.New().String(),
		name:        "Dagger",
		tags:        tag.NewContainer(tags.Melee, tags.SimpleWeapon),
		damageTimes: 1,
		damageSides: 4,
		damageTags:  tag.NewContainer(tags.Piercing),
		reach:       1,
	}
}

func NewGreatAxe() *Weapon {
	return &Weapon{
		archetype:   "great-axe",
		id:          uuid.New().String(),
		name:        "Great Axe",
		tags:        tag.NewContainer(tags.Melee, tags.MartialAxe),
		damageTimes: 2,
		damageSides: 6,
		damageTags:  tag.NewContainer(tags.Slashing),
		reach:       1,
	}
}
