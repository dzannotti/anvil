package event

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event/snapshot"
)

type UseAction struct {
	Source snapshot.Creature
	Target snapshot.Creature
	Action snapshot.Action
}

func NewUseAction(action definition.Action, source definition.Creature, target definition.Creature) (string, UseAction) {
	return "use_action", UseAction{Action: snapshot.CaptureAction(action), Source: snapshot.CaptureCreature(source), Target: snapshot.CaptureCreature(target)}
}
