package core

import "anvil/internal/expression"

func (a *Actor) ArmorClass() expression.Expression {
	expression := expression.FromScalar(10, "Base")
	expression.Evaluate()
	return expression
}
