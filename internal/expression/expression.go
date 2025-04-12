package expression

type Expression struct {
	Terms []Term
	Value int
	rng   DiceRoller
}

func (e *Expression) Clone() Expression {
	terms := make([]Term, len(e.Terms))
	for i := range e.Terms {
		terms[i] = e.Terms[i].Clone()
	}
	return Expression{
		Value: e.Value,
		Terms: terms,
	}
}
