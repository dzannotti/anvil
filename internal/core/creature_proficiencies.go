package core

import (
	"anvil/internal/tag"
)

type Proficiencies struct {
	skills tag.Container
	bonus  int
}

func NewProficiencies(bonus int) Proficiencies {
	return Proficiencies{
		skills: tag.NewContainer(),
		bonus:  bonus,
	}
}

func (p *Proficiencies) Add(tag tag.Tag) {
	p.skills.AddTag(tag)
}

func (p Proficiencies) Has(tags tag.Container) bool {
	return tags.MatchAnyTag(p.skills)
}

func (p Proficiencies) Value(tags tag.Container) int {
	if p.Has(tags) {
		return p.bonus
	}
	return 0
}
