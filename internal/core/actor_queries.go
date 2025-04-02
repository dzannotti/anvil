package core

func (a Actor) Enemies() []*Actor {
	opponents := TeamEnemies
	if a.Team == TeamEnemies {
		opponents = TeamPlayers
	}
	enemies := make([]*Actor, 0)
	for _, c := range a.Encounter.Actors {
		if opponents == c.Team {
			enemies = append(enemies, c)
		}
	}
	return enemies
}

func (a Actor) HitPointsNormalized() float32 {
	return float32(a.HitPoints) / float32(a.MaxHitPoints)
}
