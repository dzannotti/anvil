package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type KeyBinds struct {
	SelectAction func(idx int)
}

func (kb *KeyBinds) Update() {
	for i := 1; i <= 9; i++ {
		if rl.IsKeyReleased(rl.KeyOne + int32(i-1)) {
			kb.SelectAction(i)
		}
	}
}
