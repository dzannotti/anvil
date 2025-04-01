package core

import (
	"anvil/internal/expression"
	"anvil/internal/tag"
)

const (
	BeforeAttackRollStateType     = "BeforeAttackRollState"
	AfterAttackRollStateType      = "AfterAttackRollState"
	AttributeCalculationStateType = "AttributeCalculationState"
)

type BeforeAttackRollState struct {
	Source     *Creature
	Target     *Creature
	Expression *expression.Expression
	Tags       tag.Container
}

type AfterAttackRollState struct {
	Source *Creature
	Target *Creature
	Result *expression.Expression
	Tags   tag.Container
}

type AttributeCalculationState struct {
	Expression *expression.Expression
	Attribute  tag.Tag
}
