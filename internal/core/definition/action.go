package definition

type Action interface {
	Name() string
	Perform(Creature)
}
