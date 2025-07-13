package loader

type AttributesDefinition struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

type ProficienciesDefinition struct {
	Skills []string
	Bonus  int
}

type ResourcesDefinition struct {
	WalkSpeed  int
	FlySpeed   int
	SwimSpeed  int
	SpellSlot1 int
	SpellSlot2 int
	SpellSlot3 int
	SpellSlot4 int
	SpellSlot5 int
	SpellSlot6 int
	SpellSlot7 int
	SpellSlot8 int
	SpellSlot9 int
}

type ActorDefinition struct {
	Name               string
	Team               string
	HitPoints          int
	MaxHitPoints       int
	SpellCastingSource string
	Attributes         AttributesDefinition
	Proficiencies      ProficienciesDefinition
	Resources          ResourcesDefinition
}

// TODO: Add YAML loading functions like LoadActorFromFile, etc.
