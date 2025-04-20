package render

import (
	"bytes"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ScrollText struct {
	Lines      []string
	Rect       Rectangle
	Scroll     float32
	LineHeight int
	Padding    int
	BgColor    rl.Color
	TextColor  rl.Color
	FontSize   int
}

func (log *ScrollText) AddLine(line string) {
	log.Lines = append(log.Lines, line)

	// Automatically scroll to the bottom
	totalHeight := len(log.Lines) * log.LineHeight
	visibleHeight := int(log.Rect.Height)
	maxScroll := float32(totalHeight - visibleHeight)
	if maxScroll < 0 {
		maxScroll = 0
	}
	log.Scroll = maxScroll
}

// Draw renders the text log and handles scrolling.
func (log *ScrollText) Draw() {
	FillRectangle(log.Rect, log.BgColor)
	DrawRectangle(log.Rect, Black, 2)

	// Scroll handling (if mouse is over)
	mouse := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mouse, log.Rect.toRaylib()) {
		delta := rl.GetMouseWheelMove()
		log.Scroll -= delta * float32(log.LineHeight)
	}

	// Clamp scroll
	maxScroll := float32(len(log.Lines)*log.LineHeight - int(log.Rect.Height))
	if maxScroll < 0 {
		maxScroll = 0
	}
	if log.Scroll < 0 {
		log.Scroll = 0
	}
	if log.Scroll > maxScroll {
		log.Scroll = maxScroll
	}

	// Start drawing lines with clipping
	rl.BeginScissorMode(int32(log.Rect.X), int32(log.Rect.Y), int32(log.Rect.Width), int32(log.Rect.Height))
	y := int(log.Rect.Y) + log.Padding - int(log.Scroll)
	for _, line := range log.Lines {
		pos := Vector2i{X: int(log.Rect.X) + log.Padding, Y: y}
		DrawText(line, pos, log.TextColor, log.FontSize)
		y += log.LineHeight
	}
	rl.EndScissorMode()

	barWidth := 8
	barX := log.Rect.X + log.Rect.Width - barWidth
	barY := log.Rect.Y
	barHeight := log.Rect.Height

	padding := 2
	scrollAreaHeight := barHeight - padding*2

	totalContentHeight := len(log.Lines) * log.LineHeight
	if totalContentHeight > int(log.Rect.Height) {
		visibleRatio := float32(log.Rect.Height) / float32(totalContentHeight)
		scrollbarHeight := int(visibleRatio * float32(scrollAreaHeight))

		// Ensure minimum size for visibility
		if scrollbarHeight < 16 {
			scrollbarHeight = 16
		}

		scrollRatio := log.Scroll / float32(totalContentHeight-int(log.Rect.Height))
		scrollbarY := int(float32(barY+padding) + scrollRatio*float32(scrollAreaHeight-scrollbarHeight))

		scrollbarRect := rl.NewRectangle(float32(barX), float32(scrollbarY), float32(barWidth), float32(scrollbarHeight))
		rl.DrawRectangleRec(scrollbarRect, rl.Gray)
		rl.DrawRectangleLinesEx(scrollbarRect, 1, rl.DarkGray)
	}
}

func (log *ScrollText) Write(p []byte) (n int, err error) {
	lines := bytes.Split(p, []byte{'\n'})
	for i, line := range lines {
		if i == len(lines)-1 && len(line) == 0 {
			continue
		}
		log.AddLine(string(line))
	}
	return len(p), nil
}
