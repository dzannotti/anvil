package expression

type mockRoller struct {
	mockReturns   []int
	mockReturnIdx int
}

func (rng *mockRoller) Roll(sides int) int {
	val := rng.mockReturns[rng.mockReturnIdx]
	rng.mockReturnIdx++
	return val
}
