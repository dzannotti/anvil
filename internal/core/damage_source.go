package core

import (
	"anvil/internal/expression"
	"anvil/internal/tag"
)

type DamageSource interface {
	Name() string
	Damage() *expression.Expression
	Tags() *tag.Container
}

type basicDamageSource struct {
	name   string
	damage expression.Expression
	tags   tag.Container
}

func (d *basicDamageSource) Name() string {
	return d.name
}

func (d *basicDamageSource) Damage() *expression.Expression {
	dmg := d.damage.Clone()
	return &dmg
}

func (d *basicDamageSource) Tags() *tag.Container {
	tags := d.tags.Clone()
	return &tags
}

func NewDamageSource(damage expression.Expression, tags tag.Container) DamageSource {
	return &basicDamageSource{
		name:   "Basic Damage",
		damage: damage,
		tags:   tags,
	}
}
