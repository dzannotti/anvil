package expression

type Expression struct {
	Value      int
	Components []Component
	Rng        Roller
}

func (e *Expression) Evaluate() *Expression {
	e.Value = 0
	ctx := &Context{Rng: e.Rng}
	for _, component := range e.Components {
		e.Value += component.Evaluate(ctx)
	}
	return e
}

func (e *Expression) Clone() *Expression {
	clone := *e
	clone.Components = make([]Component, len(e.Components))
	for i, component := range e.Components {
		clone.Components[i] = component.Clone()
	}
	return &clone
}
