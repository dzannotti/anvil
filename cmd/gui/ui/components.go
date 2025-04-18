package ui

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawButton(rect Rectangle, text string, align TextAlignment, fontSize int, onClick func(), enabled bool) {
	DrawToggleButton(rect, text, align, fontSize, onClick, enabled, false)
}

func DrawToggleButton(rect Rectangle, text string, align TextAlignment, fontSize int, onClick func(), enabled bool, selected bool) {
	mo := rect.IsMouseOver()
	FillRectangle(rect, colorButtonBackground)
	if mo || selected {
		FillRectangle(rect, colorButtonHover)
		if enabled && rl.IsMouseButtonDown(rl.MouseButtonLeft) || selected {
			FillRectangle(rect, colorButtonDepressed)
		}
	}
	DrawRectangle(rect, colorButtonBorder, 2)
	DrawString(text, rect, Black, fontSize, AlignMiddle)
	if mo {
		if enabled {
			whenClicked = onClick
		} else {
			whenClicked = nil
		}
	}
}
