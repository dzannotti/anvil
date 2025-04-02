package expression

import "slices"

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
	set := make([]string, 0)
	for _, term := range e.Terms {
		if slices.Contains(set, term.Tags.ID()) {
			continue
		}
		set = append(set, term.Tags.ID())
	}
	return set
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
