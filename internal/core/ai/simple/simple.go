package simple

import (
	"errors"

	"anvil/internal/core"
	"anvil/internal/core/definition"
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
		ai.owner.Actions()[0].Perform(target)
	}
}

func (ai AI) ChooseTarget() (definition.Creature, error) {
	enemies := ai.Enemies()
	for i := range enemies {
		if !enemies[i].IsDead() {
			return enemies[i], nil
		}
	}
	return nil, errors.New("no target found")
}

func (ai AI) Enemies() []definition.Creature {
	_, enemies := ai.Teams()
	return enemies.Members()
}

func (ai AI) Teams() (definition.Team, definition.Team) {
	teams := ai.encounter.Teams()
	if teams[0].Contains(ai.owner) {
		return teams[0], teams[1]
	}
	return teams[1], teams[0]
}
