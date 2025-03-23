package event

import "anvil/internal/core/team"

type Creature struct {
	Name         string
	Team         team.Team
	HitPoints    int
	MaxHitPoints int
}

type Action struct {
	Name string
}

type UseAction struct {
	Action Action
	Source Creature
	Target Creature
}

type Death struct {
	Creature Creature
}

type Encounter struct {
	Creatures []Creature
}

type Round struct {
	Round     int
	Creatures []Creature
}

type TakeDamage struct {
	Creature Creature
	Damage   int
}

type Turn struct {
	Turn     int
	Creature Creature
}
