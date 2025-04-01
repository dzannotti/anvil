package simple

import (
	"errors"

	"anvil/internal/core"
)

type AI struct {
	encounter *core.Encounter
	owner     *core.Creature
}

func New(encounter *core.Encounter, owner *core.Creature) *AI {
	return &AI{
		encounter: encounter,
		owner:     owner,
	}
}

func (ai *AI) Play() {
	if target, err := ai.ChooseTarget(); err == nil {
		ai.owner.Actions[0].Perform(target)
	}
}

func (ai AI) ChooseTarget() (*core.Creature, error) {
	enemies := ai.Enemies()
	for i := range enemies {
		if !enemies[i].IsDead() {
			return enemies[i], nil
		}
	}
	return nil, errors.New("no target found")
}

func (ai AI) Enemies() []*core.Creature {
	team := core.TeamPlayers
	if ai.owner.Team == core.TeamPlayers {
		team = core.TeamEnemies
	}
	enemies := make([]*core.Creature, 0)
	for _, c := range ai.encounter.Creatures {
		if c.Team == team {
			enemies = append(enemies, c)
		}
	}
	return enemies
}
