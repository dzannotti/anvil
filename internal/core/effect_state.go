package core

import (
	"anvil/internal/expression"
	"anvil/internal/tag"
)

const (
	BeforeAttackRoll     = "BeforeAttackRoll"
	AfterAttackRoll      = "AfterAttackRoll"
	AttributeCalculation = "AttributeCalculation"
	BeforeTakeDamage     = "BeforeTakeDamage"
	AfterTakeDamage      = "AfterTakeDamage"
	BeforeDamageRoll     = "BeforeDamageRoll"
	AfterDamageRoll      = "AfterDamageRoll"
	BeforeSavingThrow    = "BeforeSavingThrow"
	AfterSavingThrow     = "AfterSavingThrow"
	AttributeChanged     = "AttributeChanged"
	ConditionAdded       = "ConditionAdded"
	ConditionRemoved     = "ConditionRemoved"
	TurnStarted          = "TurnStarted"
	TurnEnded            = "TurnEnded"
	Serialize            = "serialize"
	Deserialize          = "deserialize"
)

type BeforeAttackRollState struct {
	Source     *Actor
	Target     *Actor
	Expression *expression.Expression
	Tags       tag.Container
}

type AfterAttackRollState struct {
	Source *Actor
	Target *Actor
	Result *expression.Expression
	Tags   tag.Container
}

type AttributeCalculationState struct {
	Source     *Actor
	Expression *expression.Expression
	Attribute  tag.Tag
}

type BeforeTakeDamageState struct {
	Expression *expression.Expression
	Source     *Actor
}

type AfterTakeDamageState struct {
	Result       *expression.Expression
	Source       *Actor
	ActualDamage int
}

type BeforeDamageRollState struct {
	Expression *expression.Expression
	Source     *Actor
	Tags       tag.Container
}

type AfterDamageRollState struct {
	Result *expression.Expression
	Source *Actor
	Tags   tag.Container
}

type BeforeSavingThrowState struct {
	Expression      *expression.Expression
	Source          *Actor
	Attribute       tag.Tag
	DifficultyClass int
}

type AfterSavingThrowState struct {
	Result          *expression.Expression
	Source          *Actor
	Attribute       tag.Tag
	DifficultyClass int
}

type AttributeChangedState struct {
	Source    *Actor
	Attribute tag.Tag
	OldValue  int
	Value     int
}

type ConditionChangedState struct {
	Source    *Actor
	Condition tag.Tag
	From      *Effect
}

type TurnState struct {
	Source *Actor
}

type SerializeState struct {
	Operation string
	State     struct {
		Kind string
		ID   string
		Data any
	}
}
