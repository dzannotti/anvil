package base

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func NewAttributeModifierEffect() *core.Effect {
	applyAttackModifier := func(src *core.Actor, e *expression.Expression, tc tag.Container) {
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

	applySpellModifier := func(src *core.Actor, e *expression.Expression) {
		attr := src.Attribute(src.SpellCastingSource)
		attrMod := stats.AttributeModifier(attr.Value)
		attrName := fmt.Sprintf("%s", tags.ToReadable(src.SpellCastingSource))
		e.AddScalar(attrMod, fmt.Sprintf("Attribute Modifier (%s)", attrName), attr.Terms...)
	}

	fx := &core.Effect{Name: "AttributeModifier", Priority: core.PriorityBase}

	fx.WithBeforeAttackRoll(func(_ *core.Effect, s *core.BeforeAttackRollState) {
		if s.Tags.HasTag(tags.Ranged) || s.Tags.HasTag(tags.Melee) {
			applyAttackModifier(s.Source, s.Expression, s.Tags)
		}
		if s.Tags.HasTag(tags.Spell) {
			applySpellModifier(s.Source, s.Expression)
		}
	})

	fx.WithBeforeDamageRoll(func(_ *core.Effect, s *core.BeforeDamageRollState) {
		if s.Tags.HasTag(tags.Ranged) || s.Tags.HasTag(tags.Melee) {
			applyAttackModifier(s.Source, s.Expression, s.Tags)
		}
		if s.Tags.HasTag(tags.Spell) {
			applySpellModifier(s.Source, s.Expression)
		}
	})

	fx.WithBeforeSavingThrow(func(_ *core.Effect, s *core.BeforeSavingThrowState) {
		if s.Attribute.MatchExact(tags.HitPoints) {
			return
		}
		attr := s.Source.Attribute(s.Attribute)
		mod := stats.AttributeModifier(attr.Value)
		s.Expression.AddScalar(mod, "Attribute Modifier ("+tags.ToReadable(s.Attribute)+")", attr.Terms...)
	})

	return fx
}
