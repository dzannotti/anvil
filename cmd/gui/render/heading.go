package render

import (
	"fmt"

	"anvil/cmd/gui/ui"
	"anvil/internal/core"
)

func DrawHeading(encounter *core.Encounter) {
	textRect := ui.Rectangle{X: 600, Y: 10, Width: 650, Height: 20}
	best := encounter.ActiveActor().BestScoredAction()
	if best == nil {
		ui.DrawString("Best Action: End Turn", textRect, ui.White, 20, ui.AlignRight)
	} else {
		ui.DrawString(fmt.Sprintf("Best Action: %s", best.Action.Name()), textRect, ui.White, 20, ui.AlignRight)
	}
	ui.DrawString(fmt.Sprintf("Round %d - Turn: %d", encounter.Round+1, encounter.Turn+1), textRect, ui.White, 20, ui.AlignLeft)
}
