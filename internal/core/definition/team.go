package definition

type Team interface {
	Name() string
	AddMember(c Creature)
	IsDead() bool
	Members() []Creature
	Contains(c Creature) bool
}
