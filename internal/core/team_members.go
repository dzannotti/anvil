package core

func (t *Team) AddMember(creature *Creature) {
	t.Members = append(t.Members, creature)
}
