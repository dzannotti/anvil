package snapshot

import "anvil/internal/core/definition"

type Action struct {
	Name string
}

func CaptureAction(action definition.Action) Action {
	return Action{Name: action.Name()}
}
