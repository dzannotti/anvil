package core

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/encounter"
	"anvil/internal/core/team"
)

type Creature = creature.Creature

var NewCreature = creature.New

type Team = team.Team

var NewTeam = team.New

type Encounter = encounter.Encounter

var NewEncounter = encounter.New

var NewAttributes = creature.NewAttributes

type Attributes = creature.Attributes
type AttributeValues = creature.AttributeValues

var NewProficiencies = creature.NewProficiencies

type Proficiencies = creature.Proficiencies
