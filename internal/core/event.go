package core

import (
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
	Creatures []*Creature
	World     World
}

type RoundEvent struct {
	Round     int
	Creatures []*Creature
}

type TurnEvent struct {
	Turn     int
	Creature Creature
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
	Creature Creature
}

type AttackRollEvent struct {
	Source Creature
	Target Creature
}

type TakeDamageEvent struct {
	Target Creature
	Damage int
}

type UseActionEvent struct {
	Source Creature
	Target Creature
	Action Action
}
