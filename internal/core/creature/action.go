package creature

import "sync"

type Action interface {
	Perform(source *Creature, target *Creature, wg *sync.WaitGroup)
	Cost() int
}
