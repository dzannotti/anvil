package ai

import (
	"anvil/internal/ai/metrics"
	"anvil/internal/core"
	"anvil/internal/grid"
	"slices"
)

type AI interface {
	Play()
}

type Score struct {
	Action       core.Action
	Position     grid.Position
	DamageDone   int
	FriendlyFire int
	Movement     int
	Total        int
}

func ScorePosition(world *core.World, actor *core.Actor, action core.Action, pos grid.Position) Score {
	dmgDone := metrics.DamageDone{}
	friendlyFire := metrics.FriendlyFire{}
	movement := metrics.Movement{}
	affected := action.AffectedPositions([]grid.Position{pos})
	score := Score{
		Action:       action,
		Position:     pos,
		DamageDone:   dmgDone.Evaluate(world, actor, action, pos, affected),
		FriendlyFire: friendlyFire.Evaluate(world, actor, action, pos, affected),
		Movement:     movement.Evaluate(world, actor, action, pos, affected),
	}
	score.Total = score.DamageDone + score.FriendlyFire + score.Movement
	score.Total = max(score.Total, 0)
	return score
}

func ScoreAction(world *core.World, actor *core.Actor, action core.Action) []Score {
	valid := action.ValidPositions(actor.Position)
	scores := make([]Score, 0, len(valid))
	for _, pos := range valid {
		scores = append(scores, ScorePosition(world, actor, action, pos))
	}
	return scores
}

func ScoreChoices(world *core.World, actor *core.Actor) []Score {
	scores := make([]Score, 0, len(actor.Actions))
	for _, a := range actor.Actions {
		scores = append(scores, ScoreAction(world, actor, a)...)
	}
	slices.SortFunc(scores, func(a Score, b Score) int { return b.Total - a.Total })
	return scores
}

func PickBestAction(world *core.World, actor *core.Actor) (core.Action, grid.Position) {
	choices := ScoreChoices(world, actor)
	if len(choices) == 0 || choices[0].Total < 1 {
		return nil, grid.Position{}
	}
	return choices[0].Action, choices[0].Position
}

func Play(state *core.GameState) {
	actor := state.Encounter.ActiveActor()
	if !actor.CanAct() {
		return
	}
	for {
		action, pos := PickBestAction(state.World, actor)
		if action == nil {
			break
		}
		action.Perform([]grid.Position{pos})
		if state.Encounter.IsOver() {
			break
		}
	}
}
