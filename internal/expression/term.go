package expression

import "anvil/internal/tagcontainer"

type TermType string

const (
	TypeScalar        TermType = "scalar"
	TypeScalarMax     TermType = "scalar-max"
	TypeScalarReplace TermType = "scalar-replace"
	TypeScalarHalve   TermType = "scalar-halve"
	TypeDamageScalar  TermType = "scalar-damage"
	TypeDice          TermType = "dice"
	TypeDice20        TermType = "dice-20"
	TypeDamageDice    TermType = "dice-damage"
	TypeDiceReplace   TermType = "dice-replace"
)

type Term struct {
	Type            TermType                  `json:"type"`
	Value           int                       `json:"value"`
	Source          string                    `json:"source"`
	Values          []int                     `json:"values"`
	Times           int                       `json:"times"`
	Sides           int                       `json:"sides"`
	HasAdvantage    []string                  `json:"hasAdvantage"`
	HasDisadvantage []string                  `json:"hasDisadvantage"`
	Tags            tagcontainer.TagContainer `json:"tags"`
	Terms           []Term                    `json:"terms"`
	IsCritical      int                       `json:"isCritical"`
}

func (t *Term) shouldModifyRoll() bool {
	return (len(t.HasAdvantage) > 0) != (len(t.HasDisadvantage) > 0)
}

func (t Term) Clone() Term {
	cloned := make([]Term, len(t.Terms))
	for i := range t.Terms {
		cloned[i] = t.Terms[i].Clone()
	}
	newTerm := NewTerm(t.Type, t.Source, cloned...)
	newTerm.Value = t.Value
	newTerm.Times = t.Times
	newTerm.Sides = t.Sides
	newTerm.Values = append(make([]int, 0), t.Values...)
	newTerm.HasAdvantage = append(make([]string, 0), t.HasAdvantage...)
	newTerm.HasDisadvantage = append(make([]string, 0), t.HasDisadvantage...)
	newTerm.Tags = t.Tags.Clone()
	return newTerm
}

func NewTerm(termType TermType, source string, terms ...Term) Term {
	if terms == nil {
		terms = []Term{}
	}
	return Term{
		Type:            termType,
		Source:          source,
		Terms:           terms,
		Values:          []int{},
		HasAdvantage:    []string{},
		HasDisadvantage: []string{},
		IsCritical:      0,
		Tags:            tagcontainer.New(),
	}
}
