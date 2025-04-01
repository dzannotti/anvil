package core

import "anvil/internal/expression"

func (c *Creature) ArmorClass() expression.Expression {
	expression := expression.FromScalar(10, "Base")
	expression.Evaluate()
	return expression
}
