package ruleset

import (
	"slices"
	"strings"

	"anvil/internal/core"
	"anvil/internal/ruleset/base"
	"anvil/internal/ruleset/fighter"
	"anvil/internal/ruleset/item/armor"
	"anvil/internal/ruleset/item/weapon"
	"anvil/internal/ruleset/monster/undead/zombie"
	"anvil/internal/ruleset/shared"
)

func CreateAction(a *core.Actor, s core.SerializedAction) core.Action {
	if s.Kind == "Move" {
		return base.NewMoveAction(a)
	}
	if s.Kind == "Fireball" {
		return shared.NewFireballAction(a)
	}
	if s.Kind == "Slam" {
		return zombie.NewSlamAction(a)
	}
	if strings.HasPrefix(s.Kind, "Attack with") {
		return nil
	}
	panic("cannot deserialize action: " + s.Kind)
}

func CreateItem(_ *core.Actor, s core.SerializedItem) core.Item {
	if s.Kind == "Great Axe" {
		return weapon.NewGreatAxe()
	}
	if s.Kind == "Dagger" {
		return weapon.NewDagger()
	}
	if s.Kind == "ChainMail" {
		return armor.NewChainMail()
	}
	panic("cannot deserialize item: " + s.Kind)
}

func CreateEffect(s core.SerializedEffect) *core.Effect {
	ignored := []string{
		"ChainMail",
	}
	if slices.Contains(ignored, s.Kind) {
		return nil
	}
	if s.Kind == "Attribute Modifier" {
		return base.NewAttributeModifierEffect()
	}
	if s.Kind == "Proficiency Modifier" {
		return base.NewProficiencyModifierEffect()
	}
	if s.Kind == "Crit" {
		return base.NewCritEffect()
	}
	if s.Kind == "Undead Fortitude" {
		return shared.NewUndeadFortitudeEffect()
	}
	if s.Kind == "Fighting Style: Defense" {
		return fighter.NewFightingStyleDefense()
	}
	if s.Kind == "Death Saving Throw" {
		return base.NewDeathSavingThrowEffect()
	}
	if s.Kind == "Death" {
		return base.NewDeathEffect()
	}
	panic("cannot deserialize effect: " + s.Kind)
}
