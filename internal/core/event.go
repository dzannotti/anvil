package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

const (
	EncounterEventType            = "encounter"
	RoundEventType                = "round"
	TurnEventType                 = "turn"
	AttributeCalculationEventType = "attributeCalculation"
	CheckResultEventType          = "checkResult"
	ExpressionResultEventType     = "expressionResult"
	DiedEventType                 = "died"
	AttackRollEventType           = "attackRoll"
	TakeDamageEventType           = "takeDamage"
	UseActionEventType            = "useAction"
)

type EncounterEvent struct {
	Teams []definition.Team
	World World
}

type RoundEvent struct {
	Round     int
	Creatures []definition.Creature
}

type TurnEvent struct {
	Turn     int
	Creature definition.Creature
}

type AttributeCalculationEvent struct {
	Attribute  tag.Tag
	Expression expression.Expression
}

type CheckResultEvent struct {
	Value    int
	Against  int
	Critical bool
	Success  bool
}

type ExpressionResultEvent struct {
	Expression expression.Expression
}

type DiedEvent struct {
	Creature definition.Creature
}

type AttackRollEvent struct {
	Source definition.Creature
	Target definition.Creature
}

type TakeDamageEvent struct {
	Target definition.Creature
	Damage int
}

type UseActionEvent struct {
	Source definition.Creature
	Target definition.Creature
	Action definition.Action
}
