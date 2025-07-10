package basic

import (
	"slices"

	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

type WeaponAttackAction struct {
	owner        *core.Actor
	archetype    string
	id           string
	name         string
	tags         tag.Container
	cost         map[tag.Tag]int
	reach        int
	damageSource core.DamageSource
}

func NewWeaponAttackAction(owner *core.Actor, name string, damageSource core.DamageSource, reach int, actionTags tag.Container, cost map[tag.Tag]int) *WeaponAttackAction {
	a := &WeaponAttackAction{
		owner:        owner,
		archetype:    "attack",
		id:           uuid.New().String(),
		name:         name,
		tags:         actionTags,
		cost:         cost,
		reach:        reach,
		damageSource: damageSource,
	}
	a.tags.Add(tag.NewContainer(tags.Attack))
	return a
}

func (a *WeaponAttackAction) Owner() *core.Actor {
	return a.owner
}

func (a *WeaponAttackAction) Archetype() string {
	return a.archetype
}

func (a *WeaponAttackAction) ID() string {
	return a.id
}

func (a *WeaponAttackAction) Name() string {
	return a.name
}

func (a *WeaponAttackAction) Reach() int {
	return a.reach
}

func (a *WeaponAttackAction) Cost() map[tag.Tag]int {
	return a.cost
}

func (a *WeaponAttackAction) CanAfford() bool {
	return a.owner.Resources.CanAfford(a.cost)
}

func (a *WeaponAttackAction) Commit() {
	if !a.CanAfford() {
		panic("Attempt to commit action without affording cost")
	}

	for tag, amount := range a.cost {
		a.owner.ConsumeResource(tag, amount)
	}
}

func (a *WeaponAttackAction) Perform(pos []grid.Position) {
	target := a.owner.World.ActorAt(pos[0])
	a.owner.Dispatcher.Begin(core.UseActionEvent{Action: a, Source: a.owner, Target: pos})
	a.owner.Dispatcher.Emit(core.TargetEvent{Target: []*core.Actor{target}})
	defer a.owner.Dispatcher.End()
	a.Commit()
	result := a.owner.AttackRoll(target, *a.Tags())
	if result.Success {
		dmg := a.owner.DamageRoll(a, result.Critical)
		target.TakeDamage(*dmg)
	}
}

func (a *WeaponAttackAction) ValidPositions(from grid.Position) []grid.Position {
	if !a.CanAfford() {
		return []grid.Position{}
	}

	shape := shapes.Circle(from, a.reach)
	valid := make([]grid.Position, 0)
	enemies := a.owner.Enemies()
	for _, pos := range shape {
		if !a.owner.World.IsValidPosition(pos) {
			continue
		}

		if pos == from {
			continue
		}

		other := a.owner.World.ActorAt(pos)
		if other == nil {
			continue
		}

		if !slices.Contains(enemies, other) {
			continue
		}

		if other.IsDead() {
			continue
		}

		valid = append(valid, pos)
	}
	return valid
}

func (a *WeaponAttackAction) AffectedPositions(tar []grid.Position) []grid.Position {
	return []grid.Position{tar[0]}
}

func (a *WeaponAttackAction) Damage() *expression.Expression {
	return a.damageSource.Damage()
}

func (a *WeaponAttackAction) Tags() *tag.Container {
	combined := tag.NewContainerFromContainer(a.tags)
	combined.Add(*a.damageSource.Tags())
	return &combined
}

func (a *WeaponAttackAction) AverageDamage() int {
	return a.Damage().ExpectedValue()
}
