package ui

import "anvil/internal/core"

func DrawRequest(world *core.World, window Window) {
	request := world.RequestManager().GetPendingRequest()
	if request == nil {
		return
	}

	FillRectangle(
		Rectangle{X: 0, Y: 0, Width: window.Width, Height: window.Height},
		Color{R: Overlay0.R, G: Overlay0.G, B: Overlay0.B, A: 220},
	)
	buttonWidth := 160
	DrawString(
		request.Target.Name,
		Rectangle{X: window.Width / 2, Y: 250, Width: 1, Height: 1},
		Crust,
		22,
		AlignMiddle,
	)
	DrawString(request.Text, Rectangle{X: window.Width / 2, Y: 300, Width: 1, Height: 1}, Crust, 22, AlignMiddle)

	buttonLeft := (window.Width - (buttonWidth+20)*len(request.Options)) / 2
	for i, option := range request.Options {
		selectOption := func() {
			request.Answer(option)
		}
		drawRequest(
			Rectangle{X: buttonLeft + i*buttonWidth + i*20, Y: 370, Width: buttonWidth, Height: 40},
			option,
			selectOption,
		)
	}
}

func drawRequest(rect Rectangle, option core.RequestOption, choose func()) {
	DrawButton(rect, option.Label, AlignMiddle, 14, choose, true)
}
