package ui

import (
	"fmt"
	"slices"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"anvil/internal/ai"
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

//nolint:cyclop // reason: cyclop here is allowed
func (am *ActionManager) Draw(cam Camera) {
	if am.Active == nil {
		return
	}
	actor := am.Encounter.ActiveActor()
	world := actor.World
	valid := am.Active.ValidPositions(actor.Position)
	choices := make(map[grid.Position]ai.Score, 0)
	for _, choice := range ai.ScoreAction(world, actor, am.Active) {
		choices[choice.Position] = choice
	}
	aiChoice, aiOk := ai.CalculateBestAIAction(world, actor)
	var best core.Action
	var bestPos grid.Position
	if aiOk {
		best = aiChoice.Action
		bestPos = aiChoice.Position
	}
	for _, pos := range valid {
		rect := RectFromPos(pos)

		FillRectangle(rect, Color{R: 223, G: 142, B: 29, A: 100})
		DrawRectangle(rect.Expand(-2, -2), Peach, 2)
		color := Text
		if best != nil && bestPos == pos && best.Name() == am.Active.Name() {
			color = Green
		}
		choice, ok := choices[pos]
		if !ok {
			choice = ai.Score{}
		}
		if ok && aiOk && choice.Action.Name() == aiChoice.Action.Name() && choice.Position == aiChoice.Position {
			choice = aiChoice
		}
		DrawString(strconv.Itoa(choice.Total), rect.Expand(0, -7), color, 14, AlignBottom)
	}
	if am.Active.Tags().MatchTag(tags.Move) {
		am.drawPath(actor, cam)
	}
	am.drawAffected(cam)
}

func (am *ActionManager) drawPath(actor *core.Actor, cam Camera) {
	worldPos := cam.GetMouseGridPosition()
	path, ok := am.World.FindPath(actor.Position, worldPos)
	if !ok {
		FillCircle(ToWorldPositionCenter(worldPos), 10, Red)
		return
	}
	positions := path.Positions()
	for i := 1; i < len(positions); i++ {
		DrawLine(ToWorldPositionCenter(positions[i-1]), ToWorldPositionCenter(positions[i]), Green, 2)
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
	go am.Active.Perform([]grid.Position{mousePos})
	am.SetActive(nil)
	if am.Encounter.IsOver() {
		am.EndTurn()
	}
	return true
}

func (am *ActionManager) drawAffected(cam Camera) {
	worldPos := cam.GetMouseGridPosition()
	affected := am.Active.AffectedPositions([]grid.Position{worldPos})
	for _, pos := range affected {
		rect := RectFromPos(pos)
		FillRectangle(rect, Color{R: 238, G: 190, B: 190, A: 100})
		DrawRectangle(rect.Expand(-2, -2), Rosewater, 2)
	}
}
