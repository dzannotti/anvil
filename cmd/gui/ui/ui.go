package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var whenClicked func() = nil
var systemFont rl.Font

func Init() {
	systemFont = rl.LoadFont("/System/Library/Fonts/Monaco.ttf")
	rl.SetTextureFilter(systemFont.Texture, rl.FilterBilinear)
}

func Close() {
	rl.UnloadFont(systemFont)
}

func ProcessInput() {
	if !rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		return
	}
	if whenClicked == nil {
		return
	}
	whenClicked()
}
