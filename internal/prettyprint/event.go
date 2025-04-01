package prettyprint

type Creature struct {
	Name         string
	HitPoints    int `json:"hit_points"`
	MaxHitPoints int `json:"max_hit_points"`
}

type Team struct {
	Name    string
	Members []Creature
}

type Cell struct {
	Walkable bool
	Occupant Creature
}

type World struct {
	Width  int
	Height int
	Cells  [][]Cell
}

type TagContainer struct {
	Tags []string
}

type Term struct {
	Type            string `json:"kind"`
	Value           int
	Source          string
	Values          []int
	Times           int
	Sides           int
	HasAdvantage    []string `json:"has_advantage"`
	HasDisadvantage []string `json:"has_disadvantage"`
	Tags            TagContainer
	Terms           []Term
	IsCritical      int `json:"is_critical"`
}

type Expression struct {
	Terms []Term
	Value int
}

type EncounterEvent struct {
	Teams []Team
	World World
}
