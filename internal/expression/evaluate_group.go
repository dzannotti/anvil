package expression

import "anvil/internal/collection"

func (e *Expression) EvaluateGroup() *Expression {
	out := Expression{rng: defaultRoller{}}
	e.Evaluate()
	ids := collection.SliceElements(e.Id())
	groups := groupTermsBy(*e, ids)
	for i, group := range groups {
		value := 0
		for _, term := range group {
			value += term.Value
		}
		out.AddDamageScalar(value, group[0].Source, *group[0].Tags, groups[i]...)
	}
	return out.Evaluate()
}

func (e Expression) Id() []string {
	ids := make([]string, len(e.Terms))
	for i, term := range e.Terms {
		ids[i] = term.Tags.Id()
	}
	return ids
}

func groupTermsBy(e Expression, ids []string) [][]Term {
	terms := make([][]Term, len(ids))
	for i, id := range ids {
		for _, term := range e.Terms {
			if term.Tags.Id() != id {
				continue
			}
			terms[i] = append(terms[i], term)
		}
	}
	return terms
}
