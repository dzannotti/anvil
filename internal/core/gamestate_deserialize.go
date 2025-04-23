package core

import (
	"encoding/json"
	"fmt"
	"io"

	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/tag"
)

type CreateAction func(*Actor, SerializedAction) Action
type CreateItem func(*Actor, SerializedItem) *Item
type CreateEffect func(*Actor, SerializedEffect) *Effect

func LoadGame(r io.Reader, hub *eventbus.Hub, createAction CreateAction, createEffect CreateEffect, createItem CreateItem) (*GameState, error) {
	decoder := json.NewDecoder(r)

	var serialized struct {
		World     *SerializedWorld
		Encounter *SerializedEncounter
	}

	if err := decoder.Decode(&serialized); err != nil {
		return nil, fmt.Errorf("failed to decode game state: %w", err)
	}
	encounter := deserializeEncounter(hub, serialized.Encounter, createAction, createEffect, createItem)
	world := deserializeWorld(serialized.World, encounter)
	encounter.World = world
	return &GameState{World: world}, nil
}

func deserializeWorld(w *SerializedWorld, e *Encounter) *World {
	world := NewWorld(w.Width, w.Height)
	for _, c := range w.Cells {
		cell := &WorldCell{
			Tile:      c.Tile,
			Occupants: make([]*Actor, 0, len(c.Occupants)),
		}
		for _, o := range c.Occupants {
			cell.Occupants = append(cell.Occupants, e.FindActor(o))
		}
	}
	return world
}

func deserializeEncounter(hub *eventbus.Hub, e *SerializedEncounter, createAction CreateAction, createEffect CreateEffect, createItem CreateItem) *Encounter {
	encounter := &Encounter{
		Round: e.Round,
		Turn:  e.Turn,
		Log:   hub,
	}

	for _, a := range e.Actors {
		actor := deserializeActor(hub, &a, createAction, createEffect, createItem)
		actor.Encounter = encounter
		encounter.Actors = append(encounter.Actors, actor)
	}

	for i := range e.InitiativeOrder {
		for _, a := range e.Actors {
			encounter.InitiativeOrder[i] = encounter.FindActor(a.ID)
		}
	}
	return encounter
}

func deserializeActor(hub *eventbus.Hub, a *SerializedActor, createAction CreateAction, createEffect CreateEffect, createItem CreateItem) *Actor {
	actor := &Actor{
		id:                 a.ID,
		Name:               a.Name,
		Team:               a.Team,
		Log:                hub,
		Position:           a.Position,
		HitPoints:          a.HitPoints,
		MaxHitPoints:       a.MaxHitPoints,
		Attributes:         a.Attributes,
		SpellCastingSource: a.SpellCastingSource,
		Proficiencies: stats.Proficiencies{
			Skills: a.Proficiencies,
			Bonus:  a.ProficiencyBonus,
		},
		Resources: Resources{
			Current: deserializeTagMap(a.Resources),
			Max:     deserializeTagMap(a.MaxResources),
		},
		Actions: make([]Action, 0),
	}
	for i := range a.Actions {
		action := createAction(actor, a.Actions[i])
		if action != nil {
			actor.Actions = append(actor.Actions, action)
		}
	}
	for i := range a.Effects {
		actor.AddEffect(createEffect(actor, a.Effects[i]))
	}
	for i, c := range a.Conditions {
		for _, id := range c {
			src := actor.Effects.Find(id)
			actor.AddCondition(tag.FromString(i), src)
		}
	}
	return actor
}

func deserializeTagMap(m map[string]int) map[tag.Tag]int {
	tagMap := make(map[tag.Tag]int)
	for k, v := range m {
		tagMap[tag.FromString(k)] = v
	}
	return tagMap
}
