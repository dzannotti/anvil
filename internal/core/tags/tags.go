package tags

import (
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"anvil/internal/tag"
)

var (
	Primary = tag.FromString("Primary")

	Actor           = tag.FromString("Actor")
	ActorDefense    = tag.FromString("Actor.Defense")
	ActorHitPoints  = tag.FromString("Actor.Defense.HitPoints")
	ActorArmorClass = tag.FromString("Actor.Defense.ArmorClass")

	Resource                = tag.FromString("Actor.Resource")
	ResourceAction          = tag.FromString("Actor.Resource.Action")
	ResourceReaction        = tag.FromString("Actor.Resource.Reaction")
	ResourceBonusAction     = tag.FromString("Actor.Resource.BonusAction")
	ResourceLegendaryAction = tag.FromString("Actor.Resource.LegendaryAction")
	ResourceSorceryPoints   = tag.FromString("Actor.Resource.SorceryPoints")
	ResourceSpeed           = tag.FromString("Actor.Resource.Speed")
	ResourceUsedSpeed       = tag.FromString("Actor.Resource.Speed.Used")
	ResourceWalkSpeed       = tag.FromString("Actor.Resource.Speed.Walk")
	ResourceFlySpeed        = tag.FromString("Actor.Resource.Speed.Fly")
	ResourceSwimSpeed       = tag.FromString("Actor.Resource.Speed.Swim")

	ResourceSpellSlot1 = tag.FromString("Actor.Resource.SpellSlot.1")
	ResourceSpellSlot2 = tag.FromString("Actor.Resource.SpellSlot.2")
	ResourceSpellSlot3 = tag.FromString("Actor.Resource.SpellSlot.3")
	ResourceSpellSlot4 = tag.FromString("Actor.Resource.SpellSlot.4")
	ResourceSpellSlot5 = tag.FromString("Actor.Resource.SpellSlot.5")
	ResourceSpellSlot6 = tag.FromString("Actor.Resource.SpellSlot.6")
	ResourceSpellSlot7 = tag.FromString("Actor.Resource.SpellSlot.7")
	ResourceSpellSlot8 = tag.FromString("Actor.Resource.SpellSlot.8")
	ResourceSpellSlot9 = tag.FromString("Actor.Resource.SpellSlot.9")

	Attribute             = tag.FromString("Actor.Attribute")
	AttributeStrength     = tag.FromString("Actor.Attribute.Strength")
	AttributeDexterity    = tag.FromString("Actor.Attribute.Dexterity")
	AttributeConstitution = tag.FromString("Actor.Attribute.Constitution")
	AttributeIntelligence = tag.FromString("Actor.Attribute.Intelligence")
	AttributeWisdom       = tag.FromString("Actor.Attribute.Wisdom")
	AttributeCharisma     = tag.FromString("Actor.Attribute.Charisma")

	ProficiencyAcrobatics     = tag.FromString("Proficiency.Acrobatics")
	ProficiencyAnimalHandling = tag.FromString("Proficiency.AnimalHandling")
	ProficiencyArcana         = tag.FromString("Proficiency.Arcana")
	ProficiencyAthletics      = tag.FromString("Proficiency.Athletics")
	ProficiencyDeception      = tag.FromString("Proficiency.Deception")
	ProficiencyHistory        = tag.FromString("Proficiency.History")
	ProficiencyInsight        = tag.FromString("Proficiency.Insight")
	ProficiencyIntimidation   = tag.FromString("Proficiency.Intimidation")
	ProficiencyInvestigation  = tag.FromString("Proficiency.Investigation")
	ProficiencyMedicine       = tag.FromString("Proficiency.Medicine")
	ProficiencyNature         = tag.FromString("Proficiency.Nature")
	ProficiencyPerception     = tag.FromString("Proficiency.Perception")
	ProficiencyStealth        = tag.FromString("Proficiency.Stealth")
	ProficiencySurvival       = tag.FromString("Proficiency.Survival")

	ProficiencySaveStrength     = tag.FromString("Proficiency.Save.Strength")
	ProficiencySaveDexterity    = tag.FromString("Proficiency.Save.Dexterity")
	ProficiencySaveConstitution = tag.FromString("Proficiency.Save.Constitution")
	ProficiencySaveIntelligence = tag.FromString("Proficiency.Save.Intelligence")
	ProficiencySaveWisdom       = tag.FromString("Proficiency.Save.Wisdom")
	ProficiencySaveCharisma     = tag.FromString("Proficiency.Save.Charisma")

	DamageKind  = tag.FromString("Damage.Kind")
	Slashing    = tag.FromString("Damage.Kind.Slashing")
	Piercing    = tag.FromString("Damage.Kind.Piercing")
	Bludgeoning = tag.FromString("Damage.Kind.Bludgeoning")
	Poison      = tag.FromString("Damage.Kind.Poison")
	Radiant     = tag.FromString("Damage.Kind.Radiant")
	Fire        = tag.FromString("Damage.Kind.Fire")
	Force       = tag.FromString("Damage.Kind.Force")

	Melee  = tag.FromString("Melee")
	Ranged = tag.FromString("Ranged")

	Attack       = tag.FromString("Attack")
	WeaponAttack = tag.FromString("Attack.Weapon")
	Spell        = tag.FromString("Attack.Spell")
	Teleport     = tag.FromString("Teleport")

	Move  = tag.FromString("Action.Move")
	Dodge = tag.FromString("Action.Dodge")
	Help  = tag.FromString("Action.Help")
	Dash  = tag.FromString("Action.Dash")

	Evocation = tag.FromString("School.Evocation")

	Item = tag.FromString("Item")

	Weapon        = tag.FromString("Item.Weapon")
	Finesse       = tag.FromString("Item.Weapon.Finesse")
	NaturalWeapon = tag.FromString("Item.Weapon.Natural")
	MartialWeapon = tag.FromString("Item.Weapon.Martial")
	MartialAxe    = tag.FromString("Item.Weapon.Martial.Axe")
	SimpleWeapon  = tag.FromString("Item.Weapon.Simple")

	NaturalArmor = tag.FromString("Item.Armor.Natural")
	LightArmor   = tag.FromString("Item.Armor.Light")
	MediumArmor  = tag.FromString("Item.Armor.Medium")
	HeavyArmor   = tag.FromString("Item.Armor.Heavy")
	Shield       = tag.FromString("Item.Armor.Shield")

	Condition     = tag.FromString("Condition")
	Poisoned      = tag.FromString("Condition.Poisoned")
	Stable        = tag.FromString("Condition.Stable")
	Incapacitated = tag.FromString("Condition.Incapacitated")
	Unconscious   = tag.FromString("Condition.Incapacitated.Unconscious")
	Dead          = tag.FromString("Condition.Incapacitated.Unconscious.Dead")
)

func ToReadable(tag tag.Tag) string {
	ignore := []string{
		"actor",
		"proficiency",
		"damage",
		"kind",
		"item",
		"attribute",
		"resource",
		"condition",
		"defense",
	}
	keep := []string{}
	for _, part := range tag.AsStrings() {
		if slices.Contains(ignore, cases.Lower(language.Und).String(part)) {
			continue
		}
		keep = append(keep, cases.Title(language.Und).String(part))
	}
	return strings.Join(keep, " ")
}

func ToReadableShort(tag tag.Tag) string {
	long := ToReadable(tag)
	parts := strings.Split(long, " ")
	if len(parts) < 2 {
		return long
	}
	return parts[1]
}
