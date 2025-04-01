package core

type CheckResult struct {
	Value    int
	Against  int
	Critical bool
	Success  bool
}
