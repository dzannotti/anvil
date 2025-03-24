package encounter

import (
	"anvil/internal/core/team"
)

type Encounter struct {
	teams []*team.Team
}

func New(teams []*team.Team) *Encounter {
	return &Encounter{
		teams: teams,
	}
}
