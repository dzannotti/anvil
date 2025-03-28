package state

import (
	"anvil/internal/expression"
	"anvil/internal/tag"
)

type AttributeCalculation struct {
	Expression expression.Expression
	Attribute  tag.Tag
}

func (s *AttributeCalculation) Type() Type {
	return AttributeCalculationType
}
