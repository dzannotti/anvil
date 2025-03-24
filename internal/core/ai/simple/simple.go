package simple

import (
	"errors"

	"anvil/internal/core"
	"anvil/internal/core/definition"
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

func (ai SimpleAI) ChooseTarget() (definition.Creature, error) {
	enemies := ai.Enemies()
	for i := range enemies {
		if !enemies[i].IsDead() {
			return enemies[i], nil
		}
	}
	return nil, errors.New("no target found")
}

func (ai SimpleAI) Enemies() []definition.Creature {
	_, enemies := ai.Teams()
	return enemies.Members()
}

func (ai SimpleAI) Teams() (definition.Team, definition.Team) {
	teams := ai.encounter.Teams()
	if teams[0].Contains(ai.owner) {
		return teams[0], teams[1]
	}
	return teams[1], teams[0]
}
