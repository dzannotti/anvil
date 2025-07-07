package expression

func (e *Expression) addD20Modifier(source string, modifiers *[]string) {
	if len(e.Components) == 0 {
		panic("cannot modify expression with no components")
	}

	if e.Components[0].Type != TypeDice20 {
		panic("cannot modify expression with non-d20 component")
	}

	*modifiers = append(*modifiers, source)
}

func (e *Expression) WithAdvantage(source string) {
	e.addD20Modifier(source, &e.Components[0].HasAdvantage)
}

func (e *Expression) WithDisadvantage(source string) {
	e.addD20Modifier(source, &e.Components[0].HasDisadvantage)
}
