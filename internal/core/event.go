package core

import (
	"anvil/internal/core/pathfinding"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

const (
	EncounterType            = "encounter"
	RoundType                = "round"
	TurnType                 = "turn"
	AttributeCalculationType = "attributeCalculation"
	CheckResultType          = "checkResult"
	ExpressionResultType     = "expressionResult"
	ConfirmType              = "confirm"
	DeathType                = "death"
	AttackRollType           = "attackRoll"
	TakeDamageType           = "takeDamage"
	UseActionType            = "useAction"
	DamageRollType           = "damageRoll"
	EffectType               = "effect"
	AttributeChangedType     = "attributeChanged"
	SavingThrowType          = "savingThrow"
	SpendResourceType        = "spendResource"
	ConditionChangedType     = "conditionChanged"
	MoveType                 = "moveEvent"
	MoveStepType             = "moveStep"
)

type EncounterEvent struct {
	Actors []*Actor
	World  *World
}

type RoundEvent struct {
	Round  int
	Actors []*Actor
}

type TurnEvent struct {
	Turn  int
	Actor *Actor
}

type AttributeCalculationEvent struct {
	Attribute  tag.Tag
	Expression *expression.Expression
}

type CheckResultEvent struct {
	Value    int
	Against  int
	Critical bool
	Success  bool
}

type ConfirmEvent struct {
	Actor   *Actor
	Confirm bool
}

type DeathEvent struct {
	Actor *Actor
}

type ExpressionResultEvent struct {
	Expression *expression.Expression
}

type AttackRollEvent struct {
	Source *Actor
	Target *Actor
}

type TakeDamageEvent struct {
	Target *Actor
	Damage *expression.Expression
}

type UseActionEvent struct {
	Source *Actor
	Target *Actor
	Action Action
}

type DamageRollEvent struct {
	Source       *Actor
	Target       *Actor
	DamageSource []DamageSource
}

type EffectEvent struct {
	Source *Actor
	Effect *Effect
}

type AttributeChangeEvent struct {
	Source    *Actor
	Attribute tag.Tag
	Reason    string
	OldValue  int
	Value     int
}

type SavingThrowEvent struct {
	Expression      *expression.Expression
	Source          *Actor
	Attribute       tag.Tag
	DifficultyClass int
}

type SpendResourceEvent struct {
	Source   *Actor
	Resource tag.Tag
	Amount   int
}

type ConditionChangedEvent struct {
	Source    *Actor
	Condition tag.Tag
	From      *Effect
	Added     bool
}

type MoveEvent struct {
	World  *World
	Source *Actor
	From   grid.Position
	To     grid.Position
	Path   *pathfinding.Result
}

type MoveStepEvent struct {
	World  *World
	Source *Actor
	From   grid.Position
	To     grid.Position
}
