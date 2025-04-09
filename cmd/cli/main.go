package main

import (
	"fmt"
	"os"
	"time"

	"anvil/internal/ai"
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/prettyprint"
	"anvil/internal/ruleset"
	"anvil/internal/ruleset/fighter"
	"anvil/internal/ruleset/item/armor"
	"anvil/internal/ruleset/item/weapon"
	"anvil/internal/ruleset/monster/undead/zombie"
	"anvil/internal/tag"
)

func setupWorld(world *core.World) {
	walls := make([]grid.Position, 0, 256)
	for x := 0; x < world.Width(); x++ {
		walls = append(walls, grid.Position{X: x, Y: 0}, grid.Position{X: x, Y: world.Height() - 1}, grid.Position{X: 0, Y: x}, grid.Position{X: world.Width() - 1, Y: x})
		if x > 0 && x < world.Width()-2 {
			walls = append(walls, grid.Position{X: world.Width() - x, Y: x})
		}
	}
	for _, p := range walls {
		cell, _ := world.At(p)
		cell.Tile = core.Wall
	}
}

func main() {
	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(os.Stdout, msg)
	})
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
	cedric := ruleset.NewPCActor(&hub, world, grid.Position{X: 3, Y: 2}, "Cedric", 12, stats.Attributes{Strength: 16, Dexterity: 13, Constitution: 14, Intelligence: 8, Wisdom: 14, Charisma: 10}, stats.Proficiencies{Bonus: 2}, cres)
	cedric.Equip(weapon.NewGreatAxe())
	cedric.Equip(armor.NewChainMail())
	cedric.AddEffect(fighter.NewFightingStyleDefense())
	cedric.AddProficiency(tags.MartialWeapon)
	mob1 := zombie.New(&hub, world, grid.Position{X: 7, Y: 6}, "Zombie 1")
	//mob2 := zombie.New(&hub, world, grid.Position{X: 7, Y: 7}, "Zombie 2")
	//mob3 := zombie.New(&hub, world, grid.Position{X: 6, Y: 6}, "Zombie 3")
	encounter := &core.Encounter{
		Log:    &hub,
		World:  world,
		Actors: []*core.Actor{ /*wizard, */ cedric, mob1 /*, mob2, mob3*/},
	}
	gameAI := map[*core.Actor]ai.AI{
		//wizard: &ai.Simple{Encounter: encounter, Owner: wizard},
		cedric: &ai.Simple{Encounter: encounter, Owner: cedric},
		mob1:   &ai.Simple{Encounter: encounter, Owner: mob1},
		/*mob2:   &ai.Simple{Encounter: encounter, Owner: mob2},
		mob3:   &ai.Simple{Encounter: encounter, Owner: mob3},*/
	}
	start := time.Now()
	winner := encounter.Play(func(active *core.Actor) {
		gameAI[active].Play()
	})
	total := time.Since(start)
	if len(winner) == 0 {
		fmt.Println("All dead")
		return
	}
	fmt.Println("Winner:", winner)
	msPerRound := float32(total.Seconds()*1000) / float32(encounter.Round+1)
	fmt.Printf("%.2fms (%d rounds %.2fms)\n", float32(total.Microseconds())/float32(1000), encounter.Round+1, msPerRound)
}
