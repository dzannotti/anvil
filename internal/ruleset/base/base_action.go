package base

import (
	"math"

	"anvil/internal/core"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action struct {
	owner     *core.Actor
	archetype string
	id        string
	name      string
	tags      tag.Container
	cost      map[tag.Tag]int
	castRange int
	reach     int
	damage    []core.DamageSource
}

func MakeAction(
	owner *core.Actor,
	archetype string,
	id string,
	name string,
	t tag.Container,
	cost map[tag.Tag]int,
	castRange int,
	reach int,
	damage []core.DamageSource,
) Action {
	return Action{
		owner:     owner,
		archetype: archetype,
		id:        id,
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

func (a Action) Archetype() string {
	return a.archetype
}

func (a Action) ID() string {
	return a.id
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

func (a Action) Perform(_ []grid.Position, _ bool) {}

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
		a.owner.ConsumeResource(tag, amount)
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
