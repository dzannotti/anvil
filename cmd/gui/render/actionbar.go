package ui

import (
	"anvil/internal/core"
)

func DrawActions(
	state *core.GameState,
	actor *core.Actor,
	selectAction func(action core.Action),
	current core.Action,
	endTurn func(),
) {
	buttonWidth := 160
	isOver := actor.Encounter.IsOver()
	isEnabled := !isOver && state.World.Request == nil
	if actor.CanAct() {
		for i, a := range actor.Actions {
			selected := false
			if current != nil && current.Name() == a.Name() {
				selected = true
			}
			drawAction(
				isEnabled,
				Rectangle{X: i*buttonWidth + 20, Y: 670, Width: buttonWidth - 10, Height: 40},
				a,
				selectAction,
				selected,
			)
		}
	}
	DrawButton(
		Rectangle{X: len(actor.Actions)*buttonWidth + 20, Y: 670, Width: buttonWidth - 10, Height: 40},
		"End Turn",
		AlignMiddle,
		14,
		func() {
			endTurn()
		},
		isEnabled,
	)
}

func drawAction(enabled bool, rect Rectangle, action core.Action, choose func(action core.Action), selected bool) {
	DrawToggleButton(rect, action.Name(), AlignMiddle, 14, func() {
		choose(action)
	}, enabled, selected)
}
