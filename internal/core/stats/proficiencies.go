package stats

import (
	"anvil/internal/loader"
	"anvil/internal/tag"
)

type Proficiencies struct {
	Skills tag.Container
	Bonus  int
}

func NewProficienciesFromDefinition(def loader.ProficienciesDefinition) Proficiencies {
	proficiencies := Proficiencies{Bonus: def.Bonus}
	for _, skill := range def.Skills {
		proficiencies.Add(tag.FromString(skill))
	}
	return proficiencies
}

func (p *Proficiencies) Add(tag tag.Tag) {
	p.Skills.AddTag(tag)
}

func (p Proficiencies) Has(tags tag.Container) bool {
	return tags.MatchAny(p.Skills)
}

func (p Proficiencies) Value(tags tag.Container) int {
	if p.Has(tags) {
		return p.Bonus
	}
	return 0
}
