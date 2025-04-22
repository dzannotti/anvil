package simple

import (
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
	for {
		a := ai.Owner.BestScoredAction()
		if a == nil {
			break
		}
		a.Action.Perform(a.Position)
		if ai.Owner.Encounter.IsOver() {
			break
		}
	}
}
