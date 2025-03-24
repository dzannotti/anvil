package main

import (
	"anvil/internal/core"
	"anvil/internal/log"
	"anvil/internal/prettyprint"
	"fmt"
	"os"
	"time"
)

func printLog(event log.Event) {
	prettyprint.Print(os.Stdout, event)
}

func main() {
	log := log.New()
	log.AddCapturer(printLog)
	players := core.NewTeam("Players")
	enemies := core.NewTeam("Enemies")
	wizard := core.NewCreature("Wizard", 22)
	fighter := core.NewCreature("Fighter", 22)
	orc := core.NewCreature("Orc", 22)
	goblin := core.NewCreature("Goblin", 22)
	players.AddMember(wizard)
	players.AddMember(fighter)
	enemies.AddMember(orc)
	enemies.AddMember(goblin)
	encounter := core.NewEncounter([]*core.Team{players, enemies})
	start := time.Now()
	winnerCh := make(chan core.Team)
	go encounter.Play(winnerCh)
	winner := <-winnerCh
	fmt.Println("Winner:", winner.Name())
	fmt.Printf("%v elapsed\n", time.Since(start))
}
