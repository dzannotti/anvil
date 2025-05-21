package core

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/tag"
)

type CreateAction func(*Actor, SerializedAction) Action
type CreateItem func(*Actor, SerializedItem) Item
type CreateEffect func(SerializedEffect) *Effect

func (gs *GameState) Load(
	r io.Reader,
	hub *eventbus.Hub,
	createAction CreateAction,
	createEffect CreateEffect,
	createItem CreateItem,
) (*GameState, error) {
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
	for _, a := range encounter.Actors {
		a.World = world
	}
	return &GameState{World: world, Encounter: encounter}, nil
}

func deserializeWorld(w *SerializedWorld, e *Encounter) *World {
	world := NewWorld(w.Width, w.Height)
	for _, c := range w.Cells {
		cell, _ := world.At(c.Position)
		cell.Tile = c.Tile
		cell.Occupants = make([]*Actor, 0, len(c.Occupants))
		for _, o := range c.Occupants {
			cell.Occupants = append(cell.Occupants, e.FindActor(o))
		}
	}
	return world
}

func deserializeEncounter(
	hub *eventbus.Hub,
	e *SerializedEncounter,
	createAction CreateAction,
	createEffect CreateEffect,
	createItem CreateItem,
) *Encounter {
	encounter := &Encounter{
		Round:           e.Round,
		Turn:            e.Turn,
		Log:             hub,
		InitiativeOrder: make([]*Actor, len(e.InitiativeOrder)),
	}

	for _, a := range e.Actors {
		actor := deserializeActor(hub, &a, createAction, createEffect, createItem)
		actor.Encounter = encounter
		encounter.Actors = append(encounter.Actors, actor)
	}

	for i := range e.InitiativeOrder {
		encounter.InitiativeOrder[i] = encounter.FindActor(e.InitiativeOrder[i])
	}
	return encounter
}

func deserializeActor(
	hub *eventbus.Hub,
	a *SerializedActor,
	createAction CreateAction,
	createEffect CreateEffect,
	createItem CreateItem,
) *Actor {
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
		fx := createEffect(a.Effects[i])
		if fx != nil {
			st := SerializeState{Operation: "deserialize"}
			st.State.Data = a.Effects[i].Data
			fx.Evaluate(Deserialize, &st)
			actor.AddEffect(fx)
		}
	}
	for _, def := range a.Equipped {
		actor.Equip(createItem(actor, def))
	}
	for t, c := range a.Conditions {
		for _, id := range c {
			src := actor.Effects.Find(id)
			if src != nil {
				actor.AddCondition(tag.FromString(t), src)
			}
			if strings.Contains(t, "Dead") {
				actor.AddCondition(tag.FromString(t), &Effect{Name: "Dead"})
			}
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
