package definition

import "anvil/internal/log"

type Creature interface {
	Name() string
	IsDead() bool
	HitPoints() int
	MaxHitPoints() int
	StartTurn()
	Actions() []Action
	Log() *log.EventLog
	TakeDamage(damage int)
}
