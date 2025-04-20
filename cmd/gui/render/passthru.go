package ui

import rl "github.com/gen2brain/raylib-go/raylib"

func GetFrameTime() float32 {
	return rl.GetFrameTime()
}

type Color = rl.Color
type Vector2 = rl.Vector2

var Rosewater = Color{R: 242, G: 213, B: 207, A: 255} // #f2d5cf
var Flamingo = Color{R: 238, G: 190, B: 190, A: 255}  // #eebebe
var Pink = Color{R: 244, G: 184, B: 228, A: 255}      // #f4b8e4
var Mauve = Color{R: 202, G: 158, B: 230, A: 255}     // #ca9ee6
var Red = Color{R: 231, G: 130, B: 132, A: 255}       // #e78284
var Maroon = Color{R: 234, G: 153, B: 156, A: 255}    // #ea999c
var Peach = Color{R: 239, G: 159, B: 118, A: 255}     // #ef9f76
var Yellow = Color{R: 229, G: 200, B: 144, A: 255}    // #e5c890
var Green = Color{R: 166, G: 209, B: 137, A: 255}     // #a6d189
var Teal = Color{R: 129, G: 200, B: 190, A: 255}      // #81c8be
var Sky = Color{R: 153, G: 193, B: 241, A: 255}       // #99c1f1
var Sapphire = Color{R: 133, G: 193, B: 220, A: 255}  // #85c1dc
var Blue = Color{R: 140, G: 170, B: 238, A: 255}      // #8caaee
var Lavender = Color{R: 186, G: 187, B: 241, A: 255}  // #babbf1
var Text = Color{R: 198, G: 208, B: 245, A: 255}      // #c6d0f5
var Subtext1 = Color{R: 181, G: 191, B: 226, A: 255}  // #b5bfe2
var Subtext0 = Color{R: 165, G: 173, B: 206, A: 255}  // #a5adce
var Overlay2 = Color{R: 148, G: 156, B: 187, A: 255}  // #949cbb
var Overlay1 = Color{R: 131, G: 139, B: 167, A: 255}  // #838ba7
var Overlay0 = Color{R: 115, G: 121, B: 148, A: 255}  // #737994
var Surface2 = Color{R: 98, G: 104, B: 128, A: 255}   // #626880
var Surface1 = Color{R: 81, G: 87, B: 109, A: 255}    // #51576d
var Surface0 = Color{R: 65, G: 69, B: 89, A: 255}     // #414559
var Base = Color{R: 48, G: 52, B: 70, A: 255}         // #303446
var Mantle = Color{R: 41, G: 44, B: 60, A: 255}       // #292c3c
var Crust = Color{R: 35, G: 38, B: 52, A: 255}        // #232634

var colorButtonDepressed = Mauve
var colorButtonHover = Lavender
var colorButtonBackground = Blue
var colorButtonBorder = Sapphire
