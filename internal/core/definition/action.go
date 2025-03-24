package definition

type Action interface {
	Perform(Creature)
}
