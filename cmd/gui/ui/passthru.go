package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Color = rl.Color
type Vector2 = rl.Vector2

var Black = rl.Black
var White = rl.White
var Red = rl.Red
var Green = rl.Green
var Blue = rl.Blue
var Yellow = rl.Yellow
var Gray = rl.Gray
var LightGray = rl.LightGray

var colorButtonDepressed = rl.Color{R: 152, G: 156, B: 154, A: 255}
var colorButtonHover = White
var colorButtonBackground = rl.Color{R: 231, G: 232, B: 216, A: 255}
var colorButtonBorder = Black
