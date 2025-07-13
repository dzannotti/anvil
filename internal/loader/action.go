package loader

type MeleeActionDefinition struct {
	Name          string
	Cost          map[string]int
	Tags          []string
	Reach         int
	DamageFormula string
	DamageType    string
}

type RangedActionDefinition struct {
	Name          string
	Cost          map[string]int
	Tags          []string
	Range         int
	DamageFormula string
	DamageType    string
}

type SpellActionDefinition struct {
	Name          string
	Cost          map[string]int
	Tags          []string
	Level         int
	School        string
	CastingTime   string
	Range         string
	Duration      string
	DamageFormula string
	DamageType    string
}
