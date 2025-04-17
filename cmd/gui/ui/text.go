package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextAlignment int

var textSpacing = float32(0.25)

const (
	AlignTop TextAlignment = iota
	AlignLeft
	AlignMiddle
	AlignMiddleLeft
	AlignMiddleRight
	AlignRight
	AlignBottom
	AlignTopLeft
	AlignTopRight
	AlignBottomLeft
	AlignBottomRight
)

func DrawText(text string, pos Vector2i, color Color, size int) {
	rl.DrawTextEx(systemFont, text, pos.toRaylib(), float32(size), textSpacing, color)
}

func DrawString(text string, rect Rectangle, color Color, size int, align TextAlignment) {
	textSizeRaw := rl.MeasureTextEx(systemFont, text, float32(size), textSpacing)
	textSize := Vector2i{X: int(textSizeRaw.X), Y: int(textSizeRaw.Y)}

	var posX int
	switch align {
	case AlignRight, AlignTopRight, AlignBottomRight, AlignMiddleRight:
		posX = rect.X + rect.Width - textSize.X
	case AlignLeft, AlignTopLeft, AlignBottomLeft, AlignMiddleLeft:
		posX = rect.X
	case AlignMiddle, AlignTop, AlignBottom:
		posX = rect.X + (rect.Width-textSize.X)/2
	}

	var posY int
	switch align {
	case AlignBottom, AlignBottomLeft, AlignBottomRight:
		posY = rect.Y + rect.Height - textSize.Y
	case AlignTop, AlignTopLeft, AlignTopRight:
		posY = rect.Y
	case AlignMiddle, AlignLeft, AlignRight, AlignMiddleLeft, AlignMiddleRight:
		posY = rect.Y + (rect.Height-textSize.Y)/2
	}

	pos := Vector2i{X: posX, Y: posY}
	rl.DrawTextEx(systemFont, text, pos.toRaylib(), float32(size), 0.25, color)
}
