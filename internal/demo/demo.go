package demo

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/mathi"
	"anvil/internal/ruleset"
	"anvil/internal/ruleset/factories"
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
	// Initialize the registry
	ruleset.InitializeDefaultRegistry()
	registry := ruleset.DefaultRegistry

	world := core.NewWorld(10, 10)
	setupWorld(world)

	cedric := setupPlayer(registry, dispatcher, world)
	mob1, mob2 := setupEnemies(registry, dispatcher, world)

	encounter := &core.Encounter{
		Dispatcher: dispatcher,
		World:      world,
		Actors:     []*core.Actor{cedric, mob1, mob2},
	}
	return &core.GameState{World: world, Encounter: encounter}
}

func setupPlayer(registry *ruleset.Registry, dispatcher *eventbus.Dispatcher, world *core.World) *core.Actor {
	cres := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed:  5,
		tags.SpellSlot3: 1,
	}}
	cedric := factories.NewPCActor(
		registry,
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

	// Equip items through registry
	cedric.Equip(registry.NewItem("greataxe", nil))
	cedric.Equip(registry.NewItem("chainmail", nil))

	// Add effects through registry
	cedric.AddEffect(registry.NewEffect("fighting-style-defense", nil))

	cedric.AddProficiency(tags.MartialWeapon)

	// Add actions through registry
	cedric.AddAction(registry.NewAction("fireball", cedric, nil))

	return cedric
}

func setupEnemies(registry *ruleset.Registry, dispatcher *eventbus.Dispatcher, world *core.World) (*core.Actor, *core.Actor) {
	mob1 := registry.NewCreature("zombie", map[string]interface{}{
		"dispatcher": dispatcher,
		"world":      world,
		"position":   grid.Position{X: 7, Y: 6},
		"name":       "Zombie 1",
	})

	mob2 := registry.NewCreature("zombie", map[string]interface{}{
		"dispatcher": dispatcher,
		"world":      world,
		"position":   grid.Position{X: 7, Y: 7},
		"name":       "Zombie 2",
	})

	return mob1, mob2
}
