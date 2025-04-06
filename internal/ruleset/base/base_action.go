package base

import (
	"anvil/internal/core"
	"anvil/internal/grid"
	"anvil/internal/tag"
	"math"
)

type Action struct {
	owner *core.Actor
	name  string
	tags  tag.Container
	cost  map[tag.Tag]int
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

func (a Action) Perform(_ []grid.Position) {}
func (a Action) ValidPositions(_ grid.Position) []grid.Position {
	return []grid.Position{}
}

func (a Action) Commit() {
	if !a.CanAfford() {
		panic("Attempt to commit action without affording cost")
	}
	for tag, amount := range a.cost {
		a.owner.Resources.Consume(tag, amount)
		a.owner.Log.Add(core.SpendResourceType, core.SpendResourceEvent{Source: a.owner, Resource: tag, Amount: amount})
	}
}

func (a Action) ScoreAt(_ grid.Position) *core.ScoredAction {
	panic("you shouldn't call base ScoreAt - we cannot score here")
}

func (a Action) AverageDamage(ds []core.DamageSource) int {
	avg := 0
	for _, d := range ds {
		roll := float64(d.Sides+1) / 2.0
		avg += int(math.Floor(float64(d.Times) * roll))
	}
	return avg
}
