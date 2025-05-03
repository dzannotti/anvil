package ai

import (
	"reflect"
	"slices"
	"strings"

	"anvil/internal/ai/metrics"
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type AI interface {
	Play()
}

type Score struct {
	Action   core.Action
	Position grid.Position
	Metrics  map[string]int
	Total    int
}

func ScorePosition(world *core.World, actor *core.Actor, action core.Action, pos grid.Position) Score {
	affected := action.AffectedPositions([]grid.Position{pos})

	score := Score{
		Action:   action,
		Position: pos,
		Metrics:  make(map[string]int),
	}

	for _, metric := range metrics.Default {
		typeName := strings.ReplaceAll(reflect.TypeOf(metric).String(), "metrics.", "")
		score.Metrics[typeName] = metric.Evaluate(world, actor, action, pos, affected)
	}

	for _, value := range score.Metrics {
		score.Total += value
	}

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

func ScorePlan(world *core.World, actor *core.Actor, move core.Action) []Score {
	scores := make([]Score, 0, len(actor.Actions))
	valid := move.ValidPositions(actor.Position)
	plan := metrics.Plan{}
	for _, pos := range valid {
		score := Score{
			Action:   move,
			Position: pos,
			Metrics:  make(map[string]int),
		}
		score.Metrics["Plan"] = plan.Evaluate(world, actor, move, pos, []grid.Position{})
		score.Total = score.Metrics["Plan"]
		scores = append(scores, score)
	}
	slices.SortFunc(scores, func(a Score, b Score) int { return b.Total - a.Total })
	return scores
}

func ScoreChoices(world *core.World, actor *core.Actor) []Score {
	scores := make([]Score, 0, len(actor.Actions))
	var move core.Action
	for _, a := range actor.Actions {
		if a.Tags().MatchTag(tags.Move) {
			move = a
		}
		scores = append(scores, ScoreAction(world, actor, a)...)
	}
	slices.SortFunc(scores, func(a Score, b Score) int { return b.Total - a.Total })
	if move == nil {
		return scores
	}
	plan := ScorePlan(world, actor, move)
	if len(plan) > 0 && plan[0].Total > scores[0].Total {
		return []Score{plan[0]}
	}
	return scores
}

func CalculateBestAIAction(world *core.World, actor *core.Actor) (Score, bool) {
	choices := ScoreChoices(world, actor)
	if len(choices) == 0 || choices[0].Total < 1 {
		return Score{}, false
	}
	return choices[0], true
}

func PickBestAction(world *core.World, actor *core.Actor) (core.Action, grid.Position) {
	choice, ok := CalculateBestAIAction(world, actor)
	if !ok {
		return nil, grid.Position{}
	}
	return choice.Action, choice.Position
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
