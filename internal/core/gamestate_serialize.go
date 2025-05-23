package core

import (
	"encoding/json"
	"io"

	"anvil/internal/core/stats"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type SerializedWorld struct {
	Width  int
	Height int
	Cells  []SerializedWorldCell
}

type SerializedWorldCell struct {
	Position  grid.Position
	Tile      TerrainType
	Occupants []string
}

type SerializedEncounter struct {
	Round           int
	Turn            int
	InitiativeOrder []string
	Actors          []SerializedActor
}

type SerializedItem struct {
	Kind string
	Data any
}

type SerializedEffect struct {
	Kind string
	ID   string
	Data any
}

type SerializedAction struct {
	Kind string
}

type SerializedActor struct {
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
	Actions            []SerializedAction
	Equipped           []SerializedItem
	Resources          map[string]int
	MaxResources       map[string]int
	Effects            []SerializedEffect
	Conditions         map[string][]string
}

func (gs *GameState) Save(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	serialized := struct {
		World     *SerializedWorld
		Encounter *SerializedEncounter
	}{
		World:     serializeWorld(gs.World),
		Encounter: serializeEncounter(gs.Encounter),
	}
	return encoder.Encode(serialized)
}

func serializeWorld(w *World) *SerializedWorld {
	sw := &SerializedWorld{
		Width:  w.Width(),
		Height: w.Height(),
		Cells:  make([]SerializedWorldCell, 0, w.Width()*w.Height()),
	}

	for x := range w.Width() {
		for y := range w.Height() {
			pos := grid.Position{X: x, Y: y}
			cell, ok := w.At(pos)
			if !ok {
				continue
			}

			sCell := SerializedWorldCell{
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

func serializeEncounter(e *Encounter) *SerializedEncounter {
	se := &SerializedEncounter{
		Round:           e.Round,
		Turn:            e.Turn,
		InitiativeOrder: make([]string, len(e.InitiativeOrder)),
		Actors:          make([]SerializedActor, len(e.Actors)),
	}
	for i, a := range e.InitiativeOrder {
		se.InitiativeOrder[i] = a.ID()
	}
	for i, a := range e.Actors {
		se.Actors[i] = SerializedActor{
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
			Actions:            make([]SerializedAction, len(a.Actions)),
			Equipped:           make([]SerializedItem, len(a.Equipped)),
			Resources:          serializeTagMap(a.Resources.Current),
			MaxResources:       serializeTagMap(a.Resources.Max),
			Effects:            make([]SerializedEffect, len(a.Effects.effects)),
			Conditions:         serializeConditions(a.Conditions),
		}
		for j, ca := range a.Actions {
			se.Actors[i].Actions[j] = serializeAction(ca)
		}
		for j, item := range a.Equipped {
			se.Actors[i].Equipped[j] = serializeItem(item)
		}
		for j, fx := range a.Effects.effects {
			se.Actors[i].Effects[j] = serializeEffect(fx)
		}
	}
	return se
}

func serializeAction(a Action) SerializedAction {
	return SerializedAction{
		Kind: a.Name(),
	}
}

func serializeItem(i Item) SerializedItem {
	return SerializedItem{
		Kind: i.Name(),
	}
}

func serializeEffect(fx *Effect) SerializedEffect {
	st := SerializeState{Operation: "serialize"}
	fx.Evaluate(Serialize, &st)
	return SerializedEffect{
		ID:   fx.ID(),
		Kind: fx.Name,
		Data: st.State.Data,
	}
}

func serializeTagMap(m map[tag.Tag]int) map[string]int {
	sm := make(map[string]int, len(m))
	for t, v := range m {
		sm[t.AsString()] = v
	}
	return sm
}

func serializeConditions(c Conditions) map[string][]string {
	sc := make(map[string][]string, len(c.Sources))
	for t, v := range c.Sources {
		sc[t.AsString()] = make([]string, len(v))
		for i, e := range v {
			sc[t.AsString()][i] = e.ID()
		}
	}
	return sc
}
