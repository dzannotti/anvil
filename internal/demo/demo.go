package demo

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/mathi"
	"anvil/internal/ruleset/actor"
	"anvil/internal/ruleset/fighter"
	"anvil/internal/ruleset/item/armor"
	"anvil/internal/ruleset/item/weapon"
	"anvil/internal/ruleset/monster/undead/zombie"
	"anvil/internal/ruleset/shared"
	"anvil/internal/tag"
)

func setupWorld(world *core.World) {
	perimeter := 2 * (world.Width() + world.Height())
	walls := make([]grid.Position, 0, perimeter+world.Height())
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
	limit := mathi.Min(world.Width()-3, world.Height()-3)
	for y := 1; y < limit; y++ {
		walls = append(walls, grid.Position{X: world.Width() - 1 - y, Y: y})
	}
	for _, p := range walls {
		cell := world.At(p)
		if cell != nil {
			cell.Tile = core.Wall
		}
	}
}

func New(dispatcher *eventbus.Dispatcher) *core.GameState {
	world := core.NewWorld(10, 10)
	setupWorld(world)

	/*wres := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed:  5,
		tags.SpellSlot1: 1,
	}}

	wizard := ruleset.NewPCActor(&hub, world, grid.Position{X: 2, Y: 2}, "Wizard", 8, stats.Attributes{Strength: 10, Dexterity: 15, Constitution: 14, Intelligence: 16, Wisdom: 12, Charisma: 8}, stats.Proficiencies{Bonus: 2}, wres)
	wizard.Equip(weapon.NewDagger())*/

	cres := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed:  5,
		tags.SpellSlot3: 1,
	}}
	cedric := actor.NewPCActor(
		dispatcher,
		world,
		grid.Position{X: 6, Y: 6},
		"Cedric",
		12,
		stats.Attributes{Strength: 16, Dexterity: 13, Constitution: 14, Intelligence: 8, Wisdom: 14, Charisma: 10},
		stats.Proficiencies{Bonus: 2},
		cres,
	)
	cedric.SpellCastingSource = tags.Intelligence
	cedric.Equip(weapon.NewGreatAxe())
	cedric.Equip(armor.NewChainMail())
	cedric.AddEffect(fighter.NewFightingStyleDefense())
	cedric.AddProficiency(tags.MartialWeapon)
	cedric.AddAction(shared.NewFireballAction(cedric))
	mob1 := zombie.New(dispatcher, world, grid.Position{X: 7, Y: 6}, "Zombie 1")
	mob2 := zombie.New(dispatcher, world, grid.Position{X: 7, Y: 7}, "Zombie 2")
	// mob3 := zombie.New(hub, world, grid.Position{X: 6, Y: 6}, "Zombie 3")
	encounter := &core.Encounter{
		Dispatcher: dispatcher,
		World:      world,
		Actors:     []*core.Actor{ /*wizard, */ cedric, mob1, mob2 /* mob3*/},
	}
	return &core.GameState{World: world, Encounter: encounter}
}
