package simple

import (
	"anvil/internal/ai/aiutils"
	"anvil/internal/core"
)

type Simple struct {
	Encounter *core.Encounter
	Owner     *core.Actor
}

func (ai *Simple) Play() {
	if !ai.Owner.CanAct() {
		return
	}
	a := aiutils.BestAIChoice(ai.Owner)
	if a != nil {
		a.Action.Perform(a.Position)
	}
}
