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
	EndTurn   func()
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
	am.drawAffected(actor, cam)
}

func (am *ActionManager) drawPath(actor *core.Actor, cam Camera) {
	worldPos := cam.GetMouseGridPosition()
	path, ok := am.World.FindPath(actor.Position, worldPos)
	if !ok {
		FillCircle(ToWorldPositionCenter(worldPos), 10, Red)
		return
	}
	for i := 1; i < len(path.Path); i++ {
		DrawLine(ToWorldPositionCenter(path.Path[i-1]), ToWorldPositionCenter(path.Path[i]), Green, 2)
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
	if am.Encounter.IsOver() {
		am.EndTurn()
	}
	return true
}

func (am *ActionManager) drawAffected(actor *core.Actor, cam Camera) {
	worldPos := cam.GetMouseGridPosition()
	affected := am.Active.AffectedPositions([]grid.Position{worldPos})
	for _, pos := range affected {
		rect := RectFromPos(pos)
		FillRectangle(rect, Color{R: 238, G: 190, B: 190, A: 100})
		DrawRectangle(rect.Expand(-2, -2), Rosewater, 2)
	}
}
