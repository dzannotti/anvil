package base

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func NewAttributeModifierEffect() *core.Effect {
	applyModifier := func(src *core.Actor, e *expression.Expression, tc tag.Container) {
		str := src.Attribute(tags.Strength)
		dex := src.Attribute(tags.Dexterity)
		strMod := stats.AttributeModifier(str.Value)
		dexMod := stats.AttributeModifier(dex.Value)
		if tc.MatchTag(tags.Finesse) || tc.MatchTag(tags.Ranged) {
			e.AddScalar(dexMod, "Attribute Modifier (Dexterity)", dex.Terms...)
			return
		}
		e.AddScalar(strMod, "Attribute Modifier (Strength)", str.Terms...)
	}

	fx := &core.Effect{Name: "AttributeModifier", Priority: core.PriorityBase}

	fx.WithBeforeAttackRoll(func(_ *core.Effect, s *core.BeforeAttackRollState) {
		applyModifier(s.Source, s.Expression, s.Tags)
	})

	fx.WithBeforeDamageRoll(func(_ *core.Effect, s *core.BeforeDamageRollState) {
		applyModifier(s.Source, s.Expression, s.Tags)
	})

	fx.WithSavingThrow(func(_ *core.Effect, s *core.SavingThrowState) {
		applyModifier(s.Source, s.Expression, tag.ContainerFromTag(s.Attribute))
	})

	return fx
}
