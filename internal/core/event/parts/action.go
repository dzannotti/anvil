package parts

type Action struct {
	Name string
}

func NewAction(name string) Action {
	return Action{Name: name}
}
