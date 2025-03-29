package event

import (
	"anvil/internal/expression"
	"anvil/internal/tag"
)

type AttributeCalculation struct {
	Attribute  tag.Tag
	Expression expression.Expression
}

func NewAttributeCalculation(attribute tag.Tag, expression expression.Expression) AttributeCalculation {
	return AttributeCalculation{Attribute: attribute, Expression: expression}
}
