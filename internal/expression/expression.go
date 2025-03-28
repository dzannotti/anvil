package expression

type Expression struct {
	Terms []Term
	Value int
	rng   DiceRoller
}
