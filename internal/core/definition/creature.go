package definition

type Creature interface {
	Name() string
	IsDead() bool
	HitPoints() int
	MaxHitPoints() int
}
