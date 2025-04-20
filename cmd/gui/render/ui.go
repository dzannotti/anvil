package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var whenClicked func()
var systemFont rl.Font

func Init() {
	systemFont = rl.LoadFont("font.ttf")
	rl.SetTextureFilter(systemFont.Texture, rl.FilterBilinear)
}

func Close() {
	rl.UnloadFont(systemFont)
}

func Update() {
	whenClicked = nil
}

func ProcessInput() bool {
	if !rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		return false
	}
	if whenClicked == nil {
		return false
	}
	whenClicked()
	return true
}
