package stats

import (
	"anvil/internal/tag"
)

type Proficiencies struct {
	Skills tag.Container
	Bonus  int
}

func (p *Proficiencies) Add(tag tag.Tag) {
	p.Skills.AddTag(tag)
}

func (p Proficiencies) Has(tags tag.Container) bool {
	return p.Skills.MatchAny(tags)
}

func (p Proficiencies) Value(tags tag.Container) int {
	if p.Has(tags) {
		return p.Bonus
	}
	return 0
}
