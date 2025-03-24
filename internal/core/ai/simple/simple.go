package simple

import (
	"errors"

	"anvil/internal/core"
)

type SimpleAI struct {
	encounter *core.Encounter
	owner     *core.Creature
}

func New(encounter *core.Encounter, owner *core.Creature) *SimpleAI {
	return &SimpleAI{
		encounter: encounter,
		owner:     owner,
	}
}

func (ai *SimpleAI) Play() {
	if target, err := ai.ChooseTarget(); err == nil {
		ai.owner.Attack(target)
	}
}

func (ai SimpleAI) ChooseTarget() (*core.Creature, error) {
	enemies := ai.Enemies()
	for i := range enemies {
		if !enemies[i].IsDead() {
			return enemies[i], nil
		}
	}
	return nil, errors.New("no target found")
}

func (ai SimpleAI) Enemies() []*core.Creature {
	_, enemies := ai.Teams()
	return enemies.Members()
}

func (ai SimpleAI) Teams() (*core.Team, *core.Team) {
	teams := ai.encounter.Teams()
	if teams[0].Contains(ai.owner) {
		return teams[0], teams[1]
	}
	return teams[1], teams[0]
}
