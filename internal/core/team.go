package core

type Team struct {
	Name    string
	Members []*Creature
}

func NewTeam(name string) *Team {
	return &Team{
		Name:    name,
		Members: []*Creature{},
	}
}
