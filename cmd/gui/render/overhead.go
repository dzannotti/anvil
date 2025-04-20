package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"anvil/internal/grid"
)

type OverheadText struct {
	Pos      Rectangle
	Message  string
	Color    rl.Color
	Age      float32
	Lifetime float32
	OffsetY  float32
}

var scrollSpeed float32 = 20.0

func (t *OverheadText) Update(dt float32) {
	t.Age += dt
	t.OffsetY += dt * scrollSpeed
	if t.OffsetY > 1 {
		t.Pos.Y--
		t.OffsetY -= 1.0
	}
}

func (t *OverheadText) IsExpired() bool {
	return t.Age >= t.Lifetime
}

func (t *OverheadText) Draw() {
	alpha := 255 - int32((t.Age/t.Lifetime)*255)
	if alpha < 0 {
		alpha = 0
	}
	col := t.Color
	col.A = uint8(alpha)
	DrawString(t.Message, t.Pos, col, 19, AlignTop)
}

type OverheadManager struct {
	Texts []OverheadText
}

func (m *OverheadManager) Add(pos grid.Position, msg string, color rl.Color) {
	m.Texts = append(m.Texts, OverheadText{
		Pos:      Rectangle{X: pos.X*CellSize + CellSize/2, Y: pos.Y * CellSize, Width: CellSize, Height: CellSize},
		Message:  msg,
		Color:    color,
		Lifetime: 1,
	})
}

func (m *OverheadManager) Update(dt float32) {
	var remaining []OverheadText
	for i := range m.Texts {
		m.Texts[i].Update(dt)
		if !m.Texts[i].IsExpired() {
			remaining = append(remaining, m.Texts[i])
		}
	}
	m.Texts = remaining
}

func (m *OverheadManager) Draw() {
	for i := range m.Texts {
		m.Texts[i].Draw()
	}
}
