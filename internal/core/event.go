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
	ConfirmEventType              = "confirm"
	DeathEventType                = "death"
	AttackRollEventType           = "attackRoll"
	TakeDamageEventType           = "takeDamage"
	UseActionEventType            = "useAction"
)

type EncounterEvent struct {
	Actors []*Actor
	World  World
}

type RoundEvent struct {
	Round  int
	Actors []*Actor
}

type TurnEvent struct {
	Turn  int
	Actor Actor
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

type ConfirmEvent struct {
	Actor   Actor
	Confirm bool
}

type DeathEvent struct {
	Actor Actor
}

type ExpressionResultEvent struct {
	Expression expression.Expression
}

type AttackRollEvent struct {
	Source Actor
	Target Actor
}

type TakeDamageEvent struct {
	Target Actor
	Damage int
}

type UseActionEvent struct {
	Source Actor
	Target Actor
	Action Action
}
