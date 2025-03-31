package event

type CheckResult struct {
	Value    int
	Against  int
	Critical bool
	Success  bool
}

func NewCheckResult(value int, against int, critical bool, success bool) (string, CheckResult) {
	return "check_result", CheckResult{Value: value, Against: against, Critical: critical, Success: success}
}
