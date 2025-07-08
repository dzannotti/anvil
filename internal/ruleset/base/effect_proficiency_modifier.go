package base

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

var saveMap = map[tag.Tag]tag.Tag{
	tags.Strength:     tags.ProficiencySaveStrength,
	tags.Dexterity:    tags.ProficiencySaveDexterity,
	tags.Constitution: tags.ProficiencySaveConstitution,
	tags.Intelligence: tags.ProficiencySaveIntelligence,
	tags.Wisdom:       tags.ProficiencySaveWisdom,
	tags.Charisma:     tags.ProficiencySaveCharisma,
}

func NewProficiencyModifierEffect() *core.Effect {
	fx := &core.Effect{Name: "Proficiency Modifier", Priority: core.PriorityBase}

	fx.On(func(s *core.BeforeAttackRollState) {
		proficiency := s.Source.Proficiency(s.Tags)
		if proficiency != 0 {
			s.Expression.AddConstant(proficiency, "Proficiency Modifier")
		}
	})

	fx.On(func(s *core.BeforeSavingThrowState) {
		t, ok := saveMap[s.Attribute]
		if !ok {
			return
		}
		proficiency := s.Source.Proficiency(tag.NewContainer(t))
		if proficiency != 0 {
			s.Expression.AddConstant(proficiency, "Proficiency Modifier")
		}
	})

	return fx
}
