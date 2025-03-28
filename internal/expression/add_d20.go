package expression

import "errors"

var (
	ErrNoTerms = errors.New("no terms to evaluate")
	ErrNo20    = errors.New("only d20 expressions can give advantage/disadvantage")
)

func (e *Expression) addD20Modifier(source string, modifiers *[]string) (bool, error) {
	if len(e.Terms) == 0 {
		return false, ErrNoTerms
	}
	if e.Terms[0].Type != TypeDice20 {
		return false, ErrNo20
	}
	*modifiers = append(*modifiers, source)
	return true, nil
}

func (e *Expression) WithAdvantage(source string) (bool, error) {
	return e.addD20Modifier(source, &e.Terms[0].HasAdvantage)
}

func (e *Expression) WithDisadvantage(source string) (bool, error) {
	return e.addD20Modifier(source, &e.Terms[0].HasDisadvantage)
}
