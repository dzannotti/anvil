package encounter

func (e Encounter) IsOver() bool {
	alive := 0
	for _, t := range e.teams {
		if !t.IsDead() {
			alive++
		}
	}
	return alive <= 1
}
