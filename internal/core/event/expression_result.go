package event

import "anvil/internal/expression"

type ExpressionResult struct {
	Expression expression.Expression
}

func NewExpressionResult(expression expression.Expression) ExpressionResult {
	return ExpressionResult{Expression: expression}
}
