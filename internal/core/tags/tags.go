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
	Defense    = tag.FromString("Actor.Defense")
	HitPoints  = tag.FromString("Actor.Defense.HitPoints")
	ArmorClass = tag.FromString("Actor.Defense.ArmorClass")

	Actor = tag.FromString("Actor")

	// Resources
	Resource        = tag.FromString("Actor.Resource")
	Action          = tag.FromString("Actor.Resource.Action")
	Reaction        = tag.FromString("Actor.Resource.Reaction")
	BonusAction     = tag.FromString("Actor.Resource.BonusAction")
	LegendaryAction = tag.FromString("Actor.Resource.LegendaryAction")
	SorceryPoints   = tag.FromString("Actor.Resource.SorceryPoints")
	Speed           = tag.FromString("Actor.Resource.Speed")
	UsedSpeed       = tag.FromString("Actor.Resource.Speed.Used")
	WalkSpeed       = tag.FromString("Actor.Resource.Speed.Walk")
	FlySpeed        = tag.FromString("Actor.Resource.Speed.Fly")
	SwimSpeed       = tag.FromString("Actor.Resource.Speed.Swim")

	// Spell Slots
	SpellSlot1 = tag.FromString("Actor.Resource.SpellSlot.1")
	SpellSlot2 = tag.FromString("Actor.Resource.SpellSlot.2")
	SpellSlot3 = tag.FromString("Actor.Resource.SpellSlot.3")
	SpellSlot4 = tag.FromString("Actor.Resource.SpellSlot.4")
	SpellSlot5 = tag.FromString("Actor.Resource.SpellSlot.5")
	SpellSlot6 = tag.FromString("Actor.Resource.SpellSlot.6")
	SpellSlot7 = tag.FromString("Actor.Resource.SpellSlot.7")
	SpellSlot8 = tag.FromString("Actor.Resource.SpellSlot.8")
	SpellSlot9 = tag.FromString("Actor.Resource.SpellSlot.9")

	// Attributes
	Attribute    = tag.FromString("Actor.Attribute")
	Strength     = tag.FromString("Actor.Attribute.Strength")
	Dexterity    = tag.FromString("Actor.Attribute.Dexterity")
	Constitution = tag.FromString("Actor.Attribute.Constitution")
	Intelligence = tag.FromString("Actor.Attribute.Intelligence")
	Wisdom       = tag.FromString("Actor.Attribute.Wisdom")
	Charisma     = tag.FromString("Actor.Attribute.Charisma")

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

	Finesse = tag.FromString("Item.Weapon.Finesse")

	// Item
	Weapon        = tag.FromString("Item.Weapon")
	NaturalWeapon = tag.FromString("Item.Weapon.Natural")
	MartialWeapon = tag.FromString("Item.Weapon.Martial")
	MartialAxe    = tag.FromString("Item.Weapon.Martial.Axe")

	NaturalArmor = tag.FromString("Item.Armor.Natural")
	LightArmor   = tag.FromString("Item.Armor.Light")
	MediumArmor  = tag.FromString("Item.Armor.Medium")
	HeavyArmor   = tag.FromString("Item.Armor.Heavy")
	Shield       = tag.FromString("Item.Armor.Shield")

	// Condition
	Condition   = tag.FromString("Condition")
	Poisoned    = tag.FromString("Condition.Poisoned")
	Dead        = tag.FromString("Condition.Dead")
	Stable      = tag.FromString("Condition.Stable")
	Unconscious = tag.FromString("Condition.Unconscious")
)

func ToReadableTag(tag tag.Tag) string {
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
	for _, part := range tag.Strings() {
		if slices.Contains(ignore, cases.Lower(language.Und).String(part)) {
			continue
		}
		keep = append(keep, cases.Title(language.Und).String(part))
	}
	return strings.Join(keep, " ")
}
