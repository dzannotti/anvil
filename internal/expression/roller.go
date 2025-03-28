package expression

import (
	"math/rand"

	"github.com/adam-lavrik/go-imath/ix"
)

type DiceRoller interface {
	Roll(sides int) int
}

type defaultRoller struct{}

func (rng defaultRoller) Roll(sides int) int {
	sign := ix.Sign(sides)
	sides = ix.Abs(sides)
	if sides == 0 {
		return 0
	}
	return (rand.Intn(sides) + 1) * sign
}
