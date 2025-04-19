package render

import (
	"fmt"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"

	"anvil/cmd/gui/ui"
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type ActionManager struct {
	Active    core.Action
	Encounter *core.Encounter
	World     *core.World
}

func (am *ActionManager) SetActive(action core.Action) {
	am.Active = action
}

func (am *ActionManager) Draw(cam Camera) {
	if am.Active == nil {
		return
	}
	actor := am.Encounter.ActiveActor()
	valid := am.Active.ValidPositions(actor.Position)
	best := actor.BestScoredAction()
	for _, pos := range valid {
		rect := ui.Rectangle{X: pos.X * cellSize, Y: pos.Y * cellSize, Width: cellSize, Height: cellSize}
		ui.FillRectangle(rect, ui.Color{R: 255, G: 255, B: 255, A: 100})
		ui.DrawRectangle(rect.Expand(-2, -2), ui.White, 2)
		score := am.Active.ScoreAt(pos)
		if score == nil {
			ui.DrawString("---", rect, ui.Black, 13, ui.AlignBottom)
			continue
		}
		color := ui.Black
		if slices.Contains(best.Position, pos) {
			color = ui.Red
		}
		ui.DrawString(fmt.Sprintf("%.3f", score.Score), rect.Expand(0, -7), color, 13, ui.AlignBottom)
	}
	if am.Active.Tags().MatchTag(tags.Move) {
		am.drawPath(actor, cam)
	}
}

func (am *ActionManager) drawPath(actor *core.Actor, cam Camera) {
	worldPos := cam.GetMouseGridPosition()
	path, ok := am.World.FindPath(actor.Position, worldPos)
	halfSize := cellSize / 2
	if !ok {
		ui.FillCircle(ui.Vector2i{X: worldPos.X*cellSize + halfSize, Y: worldPos.Y*cellSize + halfSize}, 10, ui.Red)
		return
	}
	for i := 1; i < len(path.Path); i++ {
		ui.DrawLine(ui.Vector2i{X: path.Path[i-1].X*cellSize + halfSize, Y: path.Path[i-1].Y*cellSize + halfSize}, ui.Vector2i{X: path.Path[i].X*cellSize + halfSize, Y: path.Path[i].Y*cellSize + halfSize}, ui.Green, 2)
	}
}

func (am *ActionManager) ProcessInput(cam Camera) bool {
	if am.Active == nil {
		return false
	}
	if !rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		return false
	}
	mousePos := cam.GetMouseGridPosition()
	actor := am.Encounter.ActiveActor()
	valid := am.Active.ValidPositions(actor.Position)
	if !slices.Contains(valid, mousePos) {
		fmt.Println("Invalid Pos")
		am.SetActive(nil)
		return false
	}
	fmt.Println("Performing action")
	am.Active.Perform([]grid.Position{mousePos})
	am.SetActive(nil)
	return true
}
