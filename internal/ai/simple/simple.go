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
	a := ai.Owner.BestScoredAction()
	if a != nil {
		a.Action.Perform(a.Position)
	}
}
