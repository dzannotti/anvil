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

type LegacyDamageSource struct {
	Times  int
	Sides  int
	Source string
	tags   tag.Container
}

func (l *LegacyDamageSource) Name() string {
	return l.Source
}

func (l *LegacyDamageSource) Damage() *expression.Expression {
	expr := expression.FromDamageDice(l.Times, l.Sides, l.Source, l.tags)
	return &expr
}

func (l *LegacyDamageSource) Tags() *tag.Container {
	return &l.tags
}

// NewLegacyDamageSource creates a legacy damage source for backward compatibility
func NewLegacyDamageSource(times, sides int, source string, tags tag.Container) *LegacyDamageSource {
	return &LegacyDamageSource{
		Times:  times,
		Sides:  sides,
		Source: source,
		tags:   tags,
	}
}
