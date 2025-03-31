package expression

import (
	mapset "github.com/deckarep/golang-set/v2"
)

func (e *Expression) EvaluateGroup() *Expression {
	out := Expression{rng: defaultRoller{}}
	e.Evaluate()
	groups := e.groupTermsBy()
	for i, group := range groups {
		value := 0
		for _, term := range group {
			value += term.Value
		}
		out.AddDamageScalar(value, group[0].Source, group[0].Tags, groups[i]...)
	}
	return out.Evaluate()
}

func (e Expression) uniqueTags() []string {
	set := mapset.NewSet[string]()
	for _, term := range e.Terms {
		set.Add(term.Tags.ID())
	}
	return set.ToSlice()
}

func (e Expression) groupTermsBy() [][]Term {
	ids := e.uniqueTags()
	terms := make([][]Term, len(ids))
	for i, id := range ids {
		for _, term := range e.Terms {
			if term.Tags.ID() != id {
				continue
			}
			terms[i] = append(terms[i], term)
		}
	}
	return terms
}
