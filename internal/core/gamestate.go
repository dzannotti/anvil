package core

import (
	"encoding/json"
	"fmt"
	"io"

	"anvil/internal/core/stats"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type GameState struct {
	World     *World
	Encounter *Encounter
}

type serializedWorld struct {
	Width  int
	Height int
	Cells  []serializedWorldCell
}

type serializedWorldCell struct {
	Position  grid.Position
	Tile      TerrainType
	Occupants []string
}

type serializedEncounter struct {
	Round           int
	Turn            int
	InitiativeOrder []string
	Actors          []serializedActor
}

type serializedItem struct {
	Kind string
	Data any
}

type serializedEffect struct {
	Kind string
	Data any
}

type serializedAction struct {
	Kind string
}

type serializedActor struct {
	ID                 string
	Name               string
	Position           grid.Position
	Team               TeamID
	HitPoints          int
	MaxHitPoints       int
	Attributes         stats.Attributes
	Proficiencies      tag.Container
	ProficiencyBonus   int
	SpellCastingSource tag.Tag
	Actions            []serializedAction
	Equipped           []serializedItem
	Resources          map[tag.Tag]int
	MaxResources       map[tag.Tag]int
	Effects            []serializedEffect
	Conditions         map[tag.Tag][]string
}

func (gs *GameState) Save(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	serialized := struct {
		World     *serializedWorld
		Encounter *serializedEncounter
	}{
		World:     serializeWorld(gs.World),
		Encounter: serializeEncounter(gs.Encounter),
	}
	return encoder.Encode(serialized)
}

func serializeWorld(w *World) *serializedWorld {
	sw := &serializedWorld{
		Width:  w.Width(),
		Height: w.Height(),
		Cells:  make([]serializedWorldCell, 0, w.Width()*w.Height()),
	}

	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			pos := grid.Position{X: x, Y: y}
			cell, ok := w.At(pos)
			if !ok {
				continue
			}

			sCell := serializedWorldCell{
				Position:  pos,
				Tile:      cell.Tile,
				Occupants: make([]string, 0, len(cell.Occupants)),
			}
			for _, o := range cell.Occupants {
				sCell.Occupants = append(sCell.Occupants, o.ID())
			}
			sw.Cells = append(sw.Cells, sCell)
		}
	}
	return sw
}

func serializeEncounter(e *Encounter) *serializedEncounter {
	se := &serializedEncounter{
		Round:           e.Round,
		Turn:            e.Turn,
		InitiativeOrder: make([]string, len(e.InitiativeOrder)),
		Actors:          make([]serializedActor, len(e.Actors)),
	}
	for i, a := range e.InitiativeOrder {
		se.InitiativeOrder[i] = a.ID()
	}
	for i, a := range e.Actors {
		se.Actors[i] = serializedActor{
			ID:                 a.ID(),
			Name:               a.Name,
			Position:           a.Position,
			Team:               a.Team,
			HitPoints:          a.HitPoints,
			MaxHitPoints:       a.MaxHitPoints,
			Attributes:         a.Attributes,
			Proficiencies:      a.Proficiencies.Skills,
			ProficiencyBonus:   a.Proficiencies.Bonus,
			SpellCastingSource: a.SpellCastingSource,
			Actions:            make([]serializedAction, len(a.Actions)),
			Equipped:           make([]serializedItem, len(a.Equipped)),
			Resources:          a.Resources.Current,
			MaxResources:       a.Resources.Max,
			Effects:            make([]serializedEffect, len(a.Effects.effects)),
			Conditions:         make(map[tag.Tag][]string, 0),
		}
		for j, ca := range a.Actions {
			se.Actors[i].Actions[j] = serializeAction(ca)
		}
		for j, item := range a.Equipped {
			se.Actors[i].Equipped[j] = serializeItem(item)
		}
	}
	return se
}

func serializeAction(a Action) serializedAction {
	return serializedAction{
		Kind: a.Name(),
	}
}

func serializeItem(i Item) serializedItem {
	return serializedItem{
		Kind: i.Name(),
	}
}

func LoadGame(r io.Reader) (*GameState, error) {
	decoder := json.NewDecoder(r)

	var serialized struct {
		World     *serializedWorld
		Encounter *serializedEncounter
	}

	if err := decoder.Decode(&serialized); err != nil {
		return nil, fmt.Errorf("failed to decode game state: %w", err)
	}

	return &GameState{}, nil
}
