package expression

import (
	"math/rand/v2"

	"anvil/internal/mathi"
)

type DiceRoller interface {
	Roll(sides int) int
}

type DefaultRoller struct{}

func (rng DefaultRoller) Roll(sides int) int {
	sign := mathi.Sign(sides)
	sides = mathi.Abs(sides)
	if sides == 0 {
		return 0
	}
	return (rand.IntN(sides) + 1) * sign
}
