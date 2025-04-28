package base

import (
	"anvil/internal/core"
	"anvil/internal/grid"
	"anvil/internal/tag"
	"math"
)

type Action struct {
	owner     *core.Actor
	name      string
	tags      tag.Container
	cost      map[tag.Tag]int
	castRange int
	reach     int
	damage    []core.DamageSource
}

func MakeAction(owner *core.Actor, name string, t tag.Container, cost map[tag.Tag]int, castRange int, reach int, damage []core.DamageSource) Action {
	return Action{
		owner:     owner,
		name:      name,
		tags:      t,
		cost:      cost,
		castRange: castRange,
		reach:     reach,
		damage:    damage,
	}
}

func (a Action) Owner() *core.Actor {
	return a.owner
}

func (a Action) Name() string {
	return a.name
}

func (a Action) Tags() *tag.Container {
	return &a.tags
}

func (a Action) Cost() map[tag.Tag]int {
	return a.cost
}

func (a Action) Reach() int {
	return a.reach
}

func (a Action) CastRange() int {
	return a.castRange
}

func (a Action) Damage() []core.DamageSource {
	return a.damage
}

func (a Action) CanAfford() bool {
	return a.owner.Resources.CanAfford(a.cost)
}

func (a Action) Perform(_ []grid.Position) {}

func (a Action) ValidPositions(_ grid.Position) []grid.Position {
	return []grid.Position{}
}

func (a Action) AffectedPositions(_ grid.Position) []grid.Position {
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

func (a Action) AverageDamage() int {
	avg := 0
	for _, d := range a.damage {
		roll := float64(d.Sides+1) / 2.0
		avg += int(math.Floor(float64(d.Times) * roll))
	}
	return avg
}
