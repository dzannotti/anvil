package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"anvil/cmd/gui/ui"
	"anvil/internal/grid"
)

type OverheadText struct {
	Pos      ui.Rectangle
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
		t.Pos.Y -= 1
		t.OffsetY -= 1.0
	}
}

func (t *OverheadText) IsExpired() bool {
	return t.Age >= t.Lifetime
}

func (t *OverheadText) Draw(camera Camera) {
	alpha := 255 - int32((t.Age/t.Lifetime)*255)
	if alpha < 0 {
		alpha = 0
	}
	col := t.Color
	col.A = uint8(alpha)
	ui.DrawString(t.Message, t.Pos, col, 18, ui.AlignTop)
}

type OverheadManager struct {
	Texts []OverheadText
}

func (m *OverheadManager) Add(pos grid.Position, msg string, color rl.Color) {
	m.Texts = append(m.Texts, OverheadText{
		Pos:      ui.Rectangle{X: pos.X*cellSize + cellSize/2, Y: pos.Y * cellSize, Width: cellSize, Height: cellSize},
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

func (m *OverheadManager) Draw(camera Camera) {
	for i := range m.Texts {
		m.Texts[i].Draw(camera)
	}
}
