package core

/*func snapshotCreature(c Creature) map[string]any {
	return map[string]any{
		"Name":         c.Name,
		"HitPoints":    c.HitPoints,
		"MaxHitPoints": c.MaxHitPoints,
	}
}

/*
func snapshotTeam(t Team) map[string]any {
	creatures := make([]map[string]any, 0, len(t.Members))
	for i := range t.Members {
		creatures = append(creatures, snapshotCreature(t.Members[i]))
	}
	return map[string]any{
		"Name":    t.Name,
		"Members": creatures,
	}
}

func snapshotWorld(w World) map[string]any {
	cells := make([][]map[string]any, w.Width())
	for x := 0; x < w.Width(); x++ {
		cells[x] = make([]map[string]any, w.Height())
		for y := 0; y < w.Height(); y++ {
			cell, _ := w.At(grid.NewPosition(x, y))
			path, _ := w.Navigation().At(grid.NewPosition(x, y))
			occupant, _ := cell.Occupant()
			cells[x][y] = map[string]any{
				"Walkable": path.IsWalkable(),
				"Occupant": snapshotCreature(occupant),
			}
		}
	}
	return map[string]any{
		"Width":  w.Width(),
		"Height": w.Height(),
		"Cells":  cells,
	}
}

func snapshotAction(a Action) map[string]any {
	return map[string]any{
		"Name": a.Name(),
	}
}
*/
