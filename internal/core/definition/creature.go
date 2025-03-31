package definition

import (
	"anvil/internal/eventbus"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

type Creature interface {
	Name() string
	IsDead() bool
	HitPoints() int
	MaxHitPoints() int
	StartTurn()
	Actions() []Action
	ArmorClass() expression.Expression
	AttackRoll(target Creature, tags tag.Container) CheckResult
	Log() *eventbus.Hub
	TakeDamage(damage int)
}
