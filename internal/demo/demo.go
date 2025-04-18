package demo

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset"
	"anvil/internal/ruleset/fighter"
	"anvil/internal/ruleset/item/armor"
	"anvil/internal/ruleset/item/weapon"
	"anvil/internal/ruleset/monster/undead/zombie"
	"anvil/internal/tag"
)

func setupWorld(world *core.World) {
	walls := make([]grid.Position, 0, 256)
	for x := range world.Width() {
		walls = append(walls,
			grid.Position{X: x, Y: 0},
			grid.Position{X: x, Y: world.Height() - 1})
	}

	for y := range world.Height() {
		walls = append(walls,
			grid.Position{X: 0, Y: y},
			grid.Position{X: world.Width() - 1, Y: y})
	}
	for y := 1; y < world.Height()-2; y++ {
		walls = append(walls, grid.Position{X: world.Width() - y, Y: y})
	}
	for _, p := range walls {
		cell, _ := world.At(p)
		cell.Tile = core.Wall
	}
}

func Create(hub *eventbus.Hub) (*core.World, *core.Encounter) {
	world := core.NewWorld(10, 10)
	setupWorld(world)

	/*wres := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed:  5,
		tags.SpellSlot1: 1,
	}}

	wizard := ruleset.NewPCActor(&hub, world, grid.Position{X: 2, Y: 2}, "Wizard", 8, stats.Attributes{Strength: 10, Dexterity: 15, Constitution: 14, Intelligence: 16, Wisdom: 12, Charisma: 8}, stats.Proficiencies{Bonus: 2}, wres)
	wizard.Equip(weapon.NewDagger())*/

	cres := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed: 5,
	}}
	cedric := ruleset.NewPCActor(hub, world, grid.Position{X: 3, Y: 2}, "Cedric", 12, stats.Attributes{Strength: 16, Dexterity: 13, Constitution: 14, Intelligence: 8, Wisdom: 14, Charisma: 10}, stats.Proficiencies{Bonus: 2}, cres)
	cedric.Equip(weapon.NewGreatAxe())
	cedric.Equip(armor.NewChainMail())
	cedric.AddEffect(fighter.NewFightingStyleDefense())
	cedric.AddProficiency(tags.MartialWeapon)
	mob1 := zombie.New(hub, world, grid.Position{X: 7, Y: 6}, "Zombie 1")
	mob2 := zombie.New(hub, world, grid.Position{X: 7, Y: 7}, "Zombie 2")
	//mob3 := zombie.New(hub, world, grid.Position{X: 6, Y: 6}, "Zombie 3")
	encounter := &core.Encounter{
		Log:    hub,
		World:  world,
		Actors: []*core.Actor{ /*wizard, */ cedric, mob1, mob2 /* mob3*/},
	}
	return world, encounter
}
