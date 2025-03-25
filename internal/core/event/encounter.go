package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type Encounter struct {
	Teams []snapshot.Team
}

func NewEncounter(f []definition.Team) Encounter {
	factions := make([]snapshot.Team, 0, len(f))
	for i := range f {
		factions = append(factions, snapshot.CaptureTeam(f[i]))
	}
	return Encounter{Teams: factions}
}
