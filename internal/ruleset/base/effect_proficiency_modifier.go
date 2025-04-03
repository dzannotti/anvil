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

func NewProficiencyModifierEffect(_ *core.Actor) *core.Effect {
	fx := &core.Effect{Name: "Proficiency Modifier", Priority: core.PriorityBase}

	fx.WithBeforeAttackRoll(func(_ *core.Effect, s *core.BeforeAttackRollState) {
		proficiency := s.Source.Proficiency(s.Tags)
		if proficiency != 0 {
			s.Expression.AddScalar(proficiency, "Proficiency Modifier")
		}
	})

	fx.WithSavingThrow(func(_ *core.Effect, s *core.SavingThrowState) {
		t, ok := saveMap[s.Attribute]
		if !ok {
			return
		}
		proficiency := s.Source.Proficiency(tag.ContainerFromTag(t))
		if proficiency != 0 {
			s.Expression.AddScalar(proficiency, "Proficiency Modifier")
		}
	})

	return fx
}
