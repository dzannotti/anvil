package core

func (t Team) IsDead() bool {
	for _, c := range t.Members {
		if !c.IsDead() {
			return false
		}
	}
	return true
}

func (t Team) Contains(c *Creature) bool {
	for _, m := range t.Members {
		if m == c {
			return true
		}
	}
	return false
}
