package expression

import "slices"

func (e *Expression) EvaluateGroup() *Expression {
	out := Expression{rng: defaultRoller{}}
	e.Evaluate()
	groups := e.groupTermsBy()
	for _, group := range groups {
		value := 0
		for _, term := range group {
			value += term.Value
		}
		out.AddDamageScalar(value, group[0].Source, group[0].Tags, group...)
	}
	return out.Evaluate()
}

func (e Expression) uniqueTags() []string {
	set := make([]string, 0)
	for _, term := range e.Terms {
		tags := e.primaryTags(term.Tags)
		if slices.Contains(set, tags.ID()) {
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
			tags := e.primaryTags(term.Tags)
			if tags.ID() != id {
				continue
			}
			terms[i] = append(terms[i], term)
		}
	}
	return terms
}
