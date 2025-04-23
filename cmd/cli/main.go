package main

import (
	"fmt"
	"os"
	"time"

	"anvil/internal/ai"
	"anvil/internal/core"
	"anvil/internal/demo"
	"anvil/internal/eventbus"
	"anvil/internal/prettyprint"
)

func main() {
	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(os.Stdout, msg)
	})
	gameState := demo.New(&hub)
	encounter := gameState.Encounter

	gameAI := map[*core.Actor]ai.AI{}
	for _, a := range encounter.Actors {
		gameAI[a] = &ai.Simple{Encounter: encounter, Owner: a}
	}

	start := time.Now()
	encounter.Start()
	for !encounter.IsOver() {
		active := encounter.ActiveActor()
		gameAI[active].Play()
		encounter.EndTurn()
		break
	}
	werr := gameState.Save(os.Stdout)
	if werr != nil {
		fmt.Println("Error writing game state:", werr)
	}
	encounter.End()
	total := time.Since(start)
	winner, _ := encounter.Winner()
	if len(winner) == 0 {
		fmt.Println("All dead")
		return
	}
	fmt.Println("Winner:", string(winner))
	msPerRound := float32(total.Seconds()*1000) / float32(encounter.Round+1)
	fmt.Printf("%.2fms (%d rounds %.2fms)\n", float32(total.Microseconds())/float32(1000), encounter.Round+1, msPerRound)
}
