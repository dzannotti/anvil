package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/parts"
)

type UseAction struct {
	Source parts.Creature
	Target parts.Creature
	Action parts.Action
}

func NewUseAction(name string, source definition.Creature, target definition.Creature) UseAction {
	return UseAction{Action: parts.NewAction(name), Source: parts.NewCreature(source), Target: parts.NewCreature(target)}
}
