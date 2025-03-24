package team

func (t Team) IsDead() bool {
	for _, c := range t.members {
		if !c.IsDead() {
			return false
		}
	}
	return true
}

func (t Team) Name() string {
	return t.name
}
