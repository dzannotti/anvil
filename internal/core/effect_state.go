package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

const (
	BeforeAttackRollStateType     = "BeforeAttackRollState"
	AfterAttackRollStateType      = "AfterAttackRollState"
	AttributeCalculationStateType = "AttributeCalculationState"
)

type BeforeAttackRollState struct {
	Source     definition.Creature
	Target     definition.Creature
	Expression *expression.Expression
	Tags       tag.Container
}

type AfterAttackRollState struct {
	Source definition.Creature
	Target definition.Creature
	Result *expression.Expression
	Tags   tag.Container
}

type AttributeCalculationState struct {
	Expression *expression.Expression
	Attribute  tag.Tag
}
