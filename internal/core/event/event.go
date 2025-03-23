package event

import "anvil/internal/core/team"

func ByFaction(creatures []Creature) map[team.Team][]Creature {
	factions := make(map[team.Team][]Creature)
	for _, creature := range creatures {
		factions[creature.FactionID] = append(factions[creature.FactionID], creature)
	}
	return factions
}

type Creature struct {
	Name      string
	FactionID team.Team
	HitPoints int
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
