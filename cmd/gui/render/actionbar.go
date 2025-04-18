package render

import (
	"anvil/cmd/gui/ui"
	"anvil/internal/core"
)

func DrawActions(actor *core.Actor, selectAction func(action core.Action), current core.Action, endTurn func()) {
	buttonWidth := 160
	for i, a := range actor.Actions {
		selected := false
		if current != nil && current.Name() == a.Name() {
			selected = true
		}

		drawAction(ui.Rectangle{X: i*buttonWidth + 20, Y: 670, Width: buttonWidth - 10, Height: 40}, a, selectAction, selected)
	}
	ui.DrawButton(ui.Rectangle{X: len(actor.Actions)*buttonWidth + 20, Y: 670, Width: buttonWidth - 10, Height: 40}, "End Turn", ui.AlignMiddle, 14, func() {
		endTurn()
	}, true)
}
func drawAction(rect ui.Rectangle, action core.Action, choose func(action core.Action), selected bool) {
	ui.DrawToggleButton(rect, action.Name(), ui.AlignMiddle, 14, func() {
		choose(action)
	}, true, selected)
}
