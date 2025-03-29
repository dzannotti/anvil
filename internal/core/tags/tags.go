package tags

import (
	"anvil/internal/tag"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// Defenses
	Defense    = tag.FromString("Creature.Defense")
	HitPoints  = tag.FromString("Creature.Defense.HitPoints")
	ArmorClass = tag.FromString("Creature.Defense.ArmorClass")

	Creature = tag.FromString("Creature")

	// Resources
	Resource        = tag.FromString("Creature.Resource")
	Action          = tag.FromString("Creature.Resource.Action")
	Reaction        = tag.FromString("Creature.Resource.Reaction")
	BonusAction     = tag.FromString("Creature.Resource.BonusAction")
	LegendaryAction = tag.FromString("Creature.Resource.LegendaryAction")
	SorceryPoints   = tag.FromString("Creature.Resource.SorceryPoints")
	Speed           = tag.FromString("Creature.Resource.Speed")
	UsedSpeed       = tag.FromString("Creature.Resource.Speed.Used")
	WalkSpeed       = tag.FromString("Creature.Resource.Speed.Walk")
	FlySpeed        = tag.FromString("Creature.Resource.Speed.Fly")
	SwimSpeed       = tag.FromString("Creature.Resource.Speed.Swim")

	// Spell Slots
	SpellSlot1 = tag.FromString("Creature.Resource.SpellSlot.1")
	SpellSlot2 = tag.FromString("Creature.Resource.SpellSlot.2")
	SpellSlot3 = tag.FromString("Creature.Resource.SpellSlot.3")
	SpellSlot4 = tag.FromString("Creature.Resource.SpellSlot.4")
	SpellSlot5 = tag.FromString("Creature.Resource.SpellSlot.5")
	SpellSlot6 = tag.FromString("Creature.Resource.SpellSlot.6")
	SpellSlot7 = tag.FromString("Creature.Resource.SpellSlot.7")
	SpellSlot8 = tag.FromString("Creature.Resource.SpellSlot.8")
	SpellSlot9 = tag.FromString("Creature.Resource.SpellSlot.9")

	// Attributes
	Attribute    = tag.FromString("Creature.Attribute")
	Strength     = tag.FromString("Creature.Attribute.Strength")
	Dexterity    = tag.FromString("Creature.Attribute.Dexterity")
	Constitution = tag.FromString("Creature.Attribute.Constitution")
	Intelligence = tag.FromString("Creature.Attribute.Intelligence")
	Wisdom       = tag.FromString("Creature.Attribute.Wisdom")
	Charisma     = tag.FromString("Creature.Attribute.Charisma")

	// Proficiencies
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

	// Damage
	DamageKind  = tag.FromString("Damage.Kind")
	Slashing    = tag.FromString("Damage.Kind.Slashing")
	Piercing    = tag.FromString("Damage.Kind.Piercing")
	Bludgeoning = tag.FromString("Damage.Kind.Bludgeoning")
	Poison      = tag.FromString("Damage.Kind.Poison")
	Radiant     = tag.FromString("Damage.Kind.Radiant")

	Melee  = tag.FromString("melee")
	Ranged = tag.FromString("ranged")

	// Item
	ItemWeapon           = tag.FromString("Item.Weapon")
	ItemWeaponNatural    = tag.FromString("Item.Weapon.Natural")
	ItemWeaponMartial    = tag.FromString("Item.Weapon.Martial")
	ItemWeaponMartialAxe = tag.FromString("Item.Weapon.Martial.Axe")

	ItemArmorNatural = tag.FromString("Item.Armor.Natural")
	ItemArmorLight   = tag.FromString("Item.Armor.Light")
	ItemArmorMedium  = tag.FromString("Item.Armor.Medium")
	ItemArmorHeavy   = tag.FromString("Item.Armor.Heavy")
	ItemArmorShield  = tag.FromString("Item.Armor.Shield")

	// Condition
	Condition   = tag.FromString("Condition")
	Poisoned    = tag.FromString("Condition.Poisoned")
	Dead        = tag.FromString("Condition.Dead")
	Stable      = tag.FromString("Condition.Stable")
	Unconscious = tag.FromString("Condition.Unconscious")
)

func ToReadableTag(tag tag.Tag) string {
	ignore := []string{
		"creature",
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
	for _, part := range tag.Strings() {
		if slices.Contains(ignore, cases.Lower(language.Und).String(part)) {
			continue
		}
		keep = append(keep, cases.Title(language.Und).String(part))
	}
	return strings.Join(keep, " ")
}
