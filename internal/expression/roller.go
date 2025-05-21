package expression

import (
	"math/rand"

	"anvil/internal/mathi"
)

type DiceRoller interface {
	Roll(sides int) int
}

type defaultRoller struct{}

func (rng defaultRoller) Roll(sides int) int {
	sign := mathi.Sign(sides)
	sides = mathi.Abs(sides)
	if sides == 0 {
		return 0
	}
	return (rand.Intn(sides) + 1) * sign
}
