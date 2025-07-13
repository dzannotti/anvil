package core

import "anvil/internal/eventbus"

type GameState struct {
	World      *World
	Encounter  *Encounter
	Dispatcher *eventbus.Dispatcher
}
