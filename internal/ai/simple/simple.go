package simple

import (
	"errors"

	"anvil/internal/core"
)

type Simple struct {
	Encounter *core.Encounter
	Owner     *core.Actor
}

func (ai *Simple) Play() {
	if target, err := ai.ChooseTarget(); err == nil {
		ai.Owner.Actions[0].Perform(target)
	}
}

func (ai Simple) ChooseTarget() (*core.Actor, error) {
	enemies := ai.Enemies()
	for i := range enemies {
		if !enemies[i].IsDead() {
			return enemies[i], nil
		}
	}
	return nil, errors.New("no target found")
}

func (ai Simple) Enemies() []*core.Actor {
	team := core.TeamPlayers
	if ai.Owner.Team == core.TeamPlayers {
		team = core.TeamEnemies
	}
	enemies := make([]*core.Actor, 0)
	for _, c := range ai.Encounter.Actors {
		if c.Team == team {
			enemies = append(enemies, c)
		}
	}
	return enemies
}
