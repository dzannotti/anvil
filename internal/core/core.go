package core

import (
	"anvil/internal/core/encounter"
	"anvil/internal/core/team"
)

type Team = team.Team

var NewTeam = team.New

type Encounter = encounter.Encounter

var NewEncounter = encounter.New
