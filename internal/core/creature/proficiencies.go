package creature

import (
	"anvil/internal/tag"
	"anvil/internal/tagcontainer"
)

type Proficiencies struct {
	skills tagcontainer.TagContainer
	bonus  int
}

func NewProficiencies(bonus int) Proficiencies {
	return Proficiencies{
		skills: tagcontainer.New(),
		bonus:  bonus,
	}
}

func (p *Proficiencies) Add(tag tag.Tag) {
	p.skills.AddTag(tag)
}

func (p Proficiencies) Has(tags tagcontainer.TagContainer) bool {
	return tags.MatchAnyTag(p.skills)
}

func (p Proficiencies) Value(tags tagcontainer.TagContainer) int {
	if p.Has(tags) {
		return p.bonus
	}
	return 0
}
