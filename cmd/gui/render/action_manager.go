package ui

import (
	"fmt"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"

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
		rect := RectFromPos(pos)

		FillRectangle(rect, Color{R: 223, G: 142, B: 29, A: 100})
		DrawRectangle(rect.Expand(-2, -2), Peach, 2)
		score := am.Active.ScoreAt(pos)
		if score == nil {
			DrawString("---", rect, Text, 13, AlignBottom)
			continue
		}
		color := Text
		if slices.Contains(best.Position, pos) {
			color = Red
		}
		DrawString(fmt.Sprintf("%.3f", score.Score), rect.Expand(0, -7), color, 13, AlignBottom)
	}
	if am.Active.Tags().MatchTag(tags.Move) {
		am.drawPath(actor, cam)
	}
}

func (am *ActionManager) drawPath(actor *core.Actor, cam Camera) {
	worldPos := cam.GetMouseGridPosition()
	path, ok := am.World.FindPath(actor.Position, worldPos)
	halfSize := CellSize / 2
	if !ok {
		FillCircle(Vector2i{X: worldPos.X*CellSize + halfSize, Y: worldPos.Y*CellSize + halfSize}, 10, Red)
		return
	}
	for i := 1; i < len(path.Path); i++ {
		DrawLine(Vector2i{X: path.Path[i-1].X*CellSize + halfSize, Y: path.Path[i-1].Y*CellSize + halfSize}, Vector2i{X: path.Path[i].X*CellSize + halfSize, Y: path.Path[i].Y*CellSize + halfSize}, Green, 2)
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
