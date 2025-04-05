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
	Critical   *bool
}

type AfterTakeDamageState struct {
	Result          *expression.Expression
	Source          *Actor
	Critical        *bool
	EffectiveDamage int
}

type BeforeDamageRollState struct {
	Expression *expression.Expression
	Source     *Actor
	Critical   *bool
	Tags       tag.Container
}

type AfterDamageRollState struct {
	Result   *expression.Expression
	Source   *Actor
	Critical *bool
	Tags     tag.Container
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
	Critical        *bool
	DifficultyClass int
}

type AttributeChangedState struct {
	Source    *Actor
	Attribute tag.Tag
	OldValue  int
	Value     int
}
