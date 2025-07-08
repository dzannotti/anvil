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

