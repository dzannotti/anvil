package definition

import "anvil/internal/expression"

type CheckResult struct {
	Value      int
	Expression expression.Expression
	Critical   bool
	Success    bool
}

func NewCheckResult(value int, expression expression.Expression, critical bool, success bool) CheckResult {
	return CheckResult{
		Value:      value,
		Expression: expression,
		Critical:   critical,
		Success:    success,
	}
}
