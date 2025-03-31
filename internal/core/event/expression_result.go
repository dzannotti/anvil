package event

import "anvil/internal/expression"

type ExpressionResult struct {
	Expression expression.Expression
}

func NewExpressionResult(expression expression.Expression) (string, ExpressionResult) {
	return "expression_result", ExpressionResult{Expression: expression}
}
