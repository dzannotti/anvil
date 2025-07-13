package loader

type AttributesDefinition struct {
	Strength     int `yaml:"strength"`
	Dexterity    int `yaml:"dexterity"`
	Constitution int `yaml:"constitution"`
	Intelligence int `yaml:"intelligence"`
	Wisdom       int `yaml:"wisdom"`
	Charisma     int `yaml:"charisma"`
}

type ProficienciesDefinition struct {
	Skills []string `yaml:"skills"`
	Bonus  int      `yaml:"bonus"`
}

type ResourcesDefinition struct {
	WalkSpeed  int `yaml:"walk_speed"`
	FlySpeed   int `yaml:"fly_speed,omitempty"`
	SwimSpeed  int `yaml:"swim_speed,omitempty"`
	SpellSlot1 int `yaml:"spell_slot_1,omitempty"`
	SpellSlot2 int `yaml:"spell_slot_2,omitempty"`
	SpellSlot3 int `yaml:"spell_slot_3,omitempty"`
	SpellSlot4 int `yaml:"spell_slot_4,omitempty"`
	SpellSlot5 int `yaml:"spell_slot_5,omitempty"`
	SpellSlot6 int `yaml:"spell_slot_6,omitempty"`
	SpellSlot7 int `yaml:"spell_slot_7,omitempty"`
	SpellSlot8 int `yaml:"spell_slot_8,omitempty"`
	SpellSlot9 int `yaml:"spell_slot_9,omitempty"`
}

type ActorDefinition struct {
	Name               string                  `yaml:"name"`
	Team               string                  `yaml:"team"`
	HitPoints          int                     `yaml:"hit_points"`
	MaxHitPoints       int                     `yaml:"max_hit_points"`
	SpellCastingSource string                  `yaml:"spell_casting_source"`
	Attributes         AttributesDefinition    `yaml:"attributes"`
	Proficiencies      ProficienciesDefinition `yaml:"proficiencies"`
	Resources          ResourcesDefinition     `yaml:"resources"`
}

type WorldDefinition struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type ActionDefinition struct {
	Name      string         `yaml:"name"`
	Archetype string         `yaml:"archetype"`
	Cost      map[string]int `yaml:"cost"`
	Tags      []string       `yaml:"tags"`
	
	MeleeConfig  *MeleeActionConfig  `yaml:"melee,omitempty"`
	RangedConfig *RangedActionConfig `yaml:"ranged,omitempty"`
	SpellConfig  *SpellActionConfig  `yaml:"spell,omitempty"`
}

type MeleeActionConfig struct {
	Reach         int    `yaml:"reach"`
	DamageFormula string `yaml:"damage_formula"`
	DamageType    string `yaml:"damage_type"`
}

type RangedActionConfig struct {
	Range         int    `yaml:"range"`
	DamageFormula string `yaml:"damage_formula"`
	DamageType    string `yaml:"damage_type"`
}

type SpellActionConfig struct {
	Level         int    `yaml:"level"`
	School        string `yaml:"school"`
	CastingTime   string `yaml:"casting_time"`
	Range         string `yaml:"range"`
	Duration      string `yaml:"duration"`
	DamageFormula string `yaml:"damage_formula,omitempty"`
	DamageType    string `yaml:"damage_type,omitempty"`
}

// TODO: Add YAML loading functions like LoadActorFromFile, etc.
