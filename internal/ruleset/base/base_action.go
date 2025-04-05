package base

import (
	"anvil/internal/core"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type ScoringFunc = func(pos grid.Position) float32

type Action struct {
	owner  *core.Actor
	name   string
	tags   tag.Container
	cost   map[tag.Tag]int
	scorer ScoringFunc
}

func (a Action) Name() string {
	return a.name
}

func (a Action) Tags() tag.Container {
	return a.tags
}

func (a Action) Cost() map[tag.Tag]int {
	return a.cost
}

func (a Action) CanAfford() bool {
	return a.owner.Resources.CanAfford(a.cost)
}

func (a Action) Commit() {
	if !a.CanAfford() {
		panic("Attempt to commit action without affording cost")
	}
	for tag, amount := range a.cost {
		a.owner.Resources.Consume(tag, amount)
	}
}

func (a *AttackAction) WithScorer(s ScoringFunc) {
	a.scorer = s
}

func (a AttackAction) AIAction(pos grid.Position) *core.AIAction {
	return &core.AIAction{
		Action:   a,
		Position: []grid.Position{pos},
		Score:    a.scorer(pos),
	}
}
