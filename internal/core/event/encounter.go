package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/parts"
)

type Encounter struct {
	Teams []parts.Team
}

func NewEncounter(f []definition.Team) Encounter {
	factions := make([]parts.Team, 0, len(f))
	for i := range f {
		factions = append(factions, parts.NewFaction(f[i]))
	}
	return Encounter{Teams: factions}
}
