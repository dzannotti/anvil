package ui

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawButton(rect Rectangle, text string, align TextAlignment, fontSize int, onClick func(), enabled bool) {
	DrawToggleButton(rect, text, align, fontSize, onClick, enabled, false)
}

func DrawToggleButton(rect Rectangle, text string, align TextAlignment, fontSize int, onClick func(), enabled bool, selected bool) {
	mo := rect.IsMouseOver()
	color := colorButtonBackground
	textColor := Crust
	if !enabled {
		color = Surface0
		textColor = Text
	}
	FillRectangle(rect, color)
	if (mo || selected) && enabled {
		FillRectangle(rect, colorButtonHover)
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) || selected {
			FillRectangle(rect, colorButtonDepressed)
		}
	}
	DrawRectangle(rect, colorButtonBorder, 2)
	DrawString(text, rect, textColor, fontSize, align)
	if mo {
		if enabled {
			whenClicked = onClick
		} else {
			whenClicked = nil
		}
	}
}
