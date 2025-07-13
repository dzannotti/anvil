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
	gameState := setupGame()
	start := time.Now()

	runEncounter(gameState)

	total := time.Since(start)
	printResults(gameState.Encounter, total)
}

func setupGame() *core.GameState {
	dispatcher := eventbus.Dispatcher{}
	dispatcher.SubscribeAll(func(msg eventbus.Event) {
		prettyprint.Print(os.Stdout, msg)
	})
	gameState := demo.New(&dispatcher)
	gameState.Dispatcher = &dispatcher
	gameState.Encounter.Start()
	startRequestHandler(gameState)
	return gameState
}

func startRequestHandler(gameState *core.GameState) {
	go func() {
		for {
			if gameState.World.RequestManager().HasPendingRequest() {
				_ = gameState.World.RequestManager().AnswerDefault()
			}
			time.Sleep(2 * time.Millisecond)
		}
	}()
}

func runEncounter(gameState *core.GameState) {
	encounter := gameState.Encounter
	turnCount := 0
	maxTurns := 50

	for !encounter.IsOver() && turnCount < maxTurns {
		actor := encounter.ActiveActor()
		weights := getActorWeights(actor)
		ai.Play(gameState, weights)
		turnCount++
	}

	encounter.End()
}

func getActorWeights(actor *core.Actor) *ai.Weights {
	switch actor.Name {
	case "Cedric":
		return ai.NewDefaultWeights()
	case "Zombie 1":
		return ai.NewBerserkerWeights()
	case "Zombie 2":
		return ai.NewDefensiveWeights()
	default:
		return ai.NewDefaultWeights()
	}
}

func printResults(encounter *core.Encounter, total time.Duration) {
	winner, _ := encounter.Winner()

	printWinner(winner)
	printPerformanceStats(encounter, total)
}

func printWinner(winner core.TeamID) {
	if len(winner) == 0 {
		fmt.Println("🏁 Result: All dead")
		return
	}
	fmt.Printf("🏆 Winner: %s\n", string(winner))
}

func printPerformanceStats(encounter *core.Encounter, total time.Duration) {
	perRound := total / time.Duration(encounter.Round+1)
	avgTurnTime := calculateAvgTurnTime(total, encounter)
	fmt.Printf("   Total Time: %s\n", formatDuration(total))
	fmt.Printf("   Rounds: %d (%s avg per round)\n", encounter.Round+1, formatDuration(perRound))
	fmt.Printf("   Turns: %d (%s avg per turn)\n", getTurnCount(encounter), formatDuration(avgTurnTime))
}

func calculateAvgTurnTime(total time.Duration, encounter *core.Encounter) time.Duration {
	turnCount := getTurnCount(encounter)
	return total / time.Duration(turnCount)
}

func getTurnCount(encounter *core.Encounter) int {
	rounds := encounter.Round + 1
	turns := encounter.Turn + 1
	return rounds*len(encounter.Actors) + turns
}

func formatDuration(d time.Duration) string {
	if d >= time.Second {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
}
