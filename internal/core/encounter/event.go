package encounter

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/event"
)

func toCreatures(creatures []*creature.Creature) []event.Creature {
	parts := make([]event.Creature, len(creatures))
	for i, c := range creatures {
		parts[i] = creature.ToEvent(c)
	}
	return parts
}

func NewEncounterEvent(creatures []*creature.Creature) *event.Encounter {
	return &event.Encounter{
		Creatures: toCreatures(creatures),
	}
}

func NewRoundEvent(round int, creatures []*creature.Creature) *event.Round {
	return &event.Round{
		Round:     round,
		Creatures: toCreatures(creatures),
	}
}

func NewTurnEvent(turn int, c *creature.Creature) *event.Turn {
	return &event.Turn{
		Turn:     turn,
		Creature: creature.ToEvent(c),
	}
}
