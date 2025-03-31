package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type Encounter struct {
	Teams []snapshot.Team
	World snapshot.World
}

func NewEncounter(f []definition.Team, world definition.World) (string, Encounter) {
	factions := make([]snapshot.Team, 0, len(f))
	for i := range f {
		factions = append(factions, snapshot.CaptureTeam(f[i]))
	}
	return "encounter", Encounter{Teams: factions, World: snapshot.CaptureWorld(world)}
}
