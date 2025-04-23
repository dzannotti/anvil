package ruleset

import (
	"strings"

	"anvil/internal/core"
	"anvil/internal/ruleset/base"
)

func CreateAction(a *core.Actor, s core.SerializedAction) core.Action {
	var action core.Action
	if s.Kind == "Move" {
		action = base.NewMoveAction(a)
	}
	if strings.HasPrefix(s.Kind, "Attack with") {
		return nil
	}
	if action != nil {
		panic("cannot deserialize action: unknown kind " + s.Kind)
	}
	return action
}

func CreateItem(a *core.Actor, s core.SerializedItem) *core.Item       { return nil }
func CreateEffect(a *core.Actor, s core.SerializedEffect) *core.Effect { return nil }
