package ui

import (
	"anvil/internal/core"
)

func DrawActions(state *core.GameState, actor *core.Actor, selectAction func(action core.Action), current core.Action, endTurn func()) {
	buttonWidth := 160
	isOver := actor.Encounter.IsOver()
	if actor.CanAct() {
		for i, a := range actor.Actions {
			selected := false
			if current != nil && current.Name() == a.Name() {
				selected = true
			}

			enabled := !isOver || state.World.Request != nil

			drawAction(enabled, Rectangle{X: i*buttonWidth + 20, Y: 670, Width: buttonWidth - 10, Height: 40}, a, selectAction, selected)
		}
	}
	DrawButton(Rectangle{X: len(actor.Actions)*buttonWidth + 20, Y: 670, Width: buttonWidth - 10, Height: 40}, "End Turn", AlignMiddle, 14, func() {
		endTurn()
	}, !isOver)
}

func drawAction(enabled bool, rect Rectangle, action core.Action, choose func(action core.Action), selected bool) {
	DrawToggleButton(rect, action.Name(), AlignMiddle, 14, func() {
		choose(action)
	}, enabled, selected)
}
