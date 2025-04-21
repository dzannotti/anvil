package ui

import (
	"fmt"

	"anvil/internal/core"
)

func DrawHeading(encounter *core.Encounter) {
	textRect := Rectangle{X: 600, Y: 10, Width: 650, Height: 20}
	best := encounter.ActiveActor().BestScoredAction()
	if best == nil {
		DrawString("Best Action: End Turn", textRect, Text, 20, AlignRight)
	} else {
		DrawString(fmt.Sprintf("Best Action: %s", best.Action.Name()), textRect, Text, 20, AlignRight)
	}
	DrawString(fmt.Sprintf("Round %d - Turn: %d", encounter.Round+1, encounter.Turn+1), textRect, Text, 20, AlignLeft)
}
