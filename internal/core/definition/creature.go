package definition

type Creature interface {
	Name() string
	IsDead() bool
	HitPoints() int
	MaxHitPoints() int
	StartTurn()
	Attack(Creature)
	TakeDamage(damage int)
}
