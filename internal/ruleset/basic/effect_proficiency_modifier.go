package basic

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

var saveMap = map[tag.Tag]tag.Tag{
	tags.AttributeStrength:     tags.ProficiencySaveStrength,
	tags.AttributeDexterity:    tags.ProficiencySaveDexterity,
	tags.AttributeConstitution: tags.ProficiencySaveConstitution,
	tags.AttributeIntelligence: tags.ProficiencySaveIntelligence,
	tags.AttributeWisdom:       tags.ProficiencySaveWisdom,
	tags.AttributeCharisma:     tags.ProficiencySaveCharisma,
}

func NewProficiencyModifierEffect() *core.Effect {
	fx := &core.Effect{Name: "Proficiency Modifier", Priority: core.PriorityBase}

	fx.On(func(s *core.PreAttackRoll) {
		proficiency := s.Source.Proficiency(s.Tags)
		if proficiency != 0 {
			s.Expression.AddConstant(proficiency, "Proficiency Modifier")
		}
	})

	fx.On(func(s *core.PreSavingThrow) {
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
