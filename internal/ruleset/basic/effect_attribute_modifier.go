package basic

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

// nolint:funlen // TODO: refactor
func NewAttributeModifierEffect() *core.Effect {
	applyAttackModifier := func(src *core.Actor, e *expression.Expression, tc tag.Container) {
		str := src.Attribute(tags.AttributeStrength)
		dex := src.Attribute(tags.AttributeDexterity)
		strMod := stats.AttributeModifier(str.Value)
		dexMod := stats.AttributeModifier(dex.Value)
		if tc.MatchTag(tags.Finesse) || tc.MatchTag(tags.Ranged) {
			e.AddConstant(dexMod, "Attribute Modifier (Dexterity)", dex.Components...)
			return
		}
		e.AddConstant(strMod, "Attribute Modifier (Strength)", str.Components...)
	}

	applySpellModifier := func(src *core.Actor, e *expression.Expression) {
		attr := src.Attribute(src.SpellCastingSource)
		attrMod := stats.AttributeModifier(attr.Value)
		attrName := tags.ToReadable(src.SpellCastingSource)
		e.AddConstant(attrMod, fmt.Sprintf("Attribute Modifier (%s)", attrName), attr.Components...)
	}

	fx := &core.Effect{Name: "Attribute Modifier", Priority: core.PriorityBase}

	fx.On(func(s *core.PreAttackRoll) {
		if s.Tags.HasTag(tags.Ranged) || s.Tags.HasTag(tags.Melee) {
			applyAttackModifier(s.Source, s.Expression, s.Tags)
		}

		if s.Tags.HasTag(tags.Spell) {
			applySpellModifier(s.Source, s.Expression)
		}
	})

	fx.On(func(s *core.PreDamageRoll) {
		if s.Tags.HasTag(tags.Ranged) || s.Tags.HasTag(tags.Melee) {
			applyAttackModifier(s.Source, s.Expression, s.Tags)
		}

		if s.Tags.HasTag(tags.Spell) {
			applySpellModifier(s.Source, s.Expression)
		}
	})

	fx.On(func(s *core.PreSavingThrow) {
		if s.Attribute.MatchExact(tags.ActorHitPoints) {
			return
		}

		attr := s.Source.Attribute(s.Attribute)
		mod := stats.AttributeModifier(attr.Value)
		s.Expression.AddConstant(mod, "Attribute Modifier ("+tags.ToReadable(s.Attribute)+")", attr.Components...)
	})

	return fx
}
