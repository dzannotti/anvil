package core

type Action interface {
	Name() string
	Perform(target *Creature)
}
