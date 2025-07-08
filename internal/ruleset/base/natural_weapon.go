package base

import (
	"anvil/internal/expression"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

type NaturalWeapon struct {
	archetype string
	id        string
	name      string
	damage    expression.Expression
	tags      tag.Container
}

func NewNaturalWeapon(
	name string,
	archetype string,
	damage expression.Expression,
	tags tag.Container,
) *NaturalWeapon {
	return &NaturalWeapon{
		archetype: archetype,
		id:        uuid.New().String(),
		name:      name,
		damage:    damage,
		tags:      tags,
	}
}

func (w *NaturalWeapon) Archetype() string {
	return w.archetype
}

func (w *NaturalWeapon) ID() string {
	return w.id
}

func (w *NaturalWeapon) Name() string {
	return w.name
}

func (w *NaturalWeapon) Damage() *expression.Expression {
	return &w.damage
}

func (w *NaturalWeapon) Tags() *tag.Container {
	return &w.tags
}
