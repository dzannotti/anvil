package creature

type Creature struct {
	name         string
	hitPoints    int
	maxHitPoints int
}

func New(name string, hitPoints int) *Creature {
	return &Creature{
		name:         name,
		hitPoints:    hitPoints,
		maxHitPoints: hitPoints,
	}
}
