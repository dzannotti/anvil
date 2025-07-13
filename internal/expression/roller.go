package expression

import (
	"math/rand/v2"
	"time"
)

type Roller interface {
	Roll(sides int) int
}

type RngRoller struct {
	rng *rand.Rand
}

func NewRngRoller() *RngRoller {
	source := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()))
	return &RngRoller{rng: rand.New(source)}
}

func (r *RngRoller) Roll(sides int) int {
	return r.rng.IntN(sides) + 1
}
