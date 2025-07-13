package loader

type DamageData struct {
	Formula string
	Kind    string
}

type WeaponDefinition struct {
	Archetype string
	Name      string
	Damage    []DamageData
	Tags      []string
	Reach     int
}
