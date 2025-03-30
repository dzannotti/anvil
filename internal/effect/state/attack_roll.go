package state

import (
	"anvil/internal/core/definition"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

type BeforeAttackRoll struct {
	Source     definition.Creature
	Target     definition.Creature
	Expression *expression.Expression
	Tags       tag.Container
}

func (s *BeforeAttackRoll) Type() Type {
	return BeforeAttackRollType
}

func NewBeforeAttackRoll(source definition.Creature, target definition.Creature, expr *expression.Expression, tags tag.Container) *BeforeAttackRoll {
	return &BeforeAttackRoll{
		Source:     source,
		Target:     target,
		Expression: expr,
		Tags:       tags,
	}
}

type AfterAttackRoll struct {
	Source definition.Creature
	Target definition.Creature
	Result *expression.Expression
	Tags   tag.Container
}

func (s *AfterAttackRoll) Type() Type {
	return AfterAttackRollType
}

func NewAfterAttackRoll(source definition.Creature, target definition.Creature, result *expression.Expression, tags tag.Container) *AfterAttackRoll {
	return &AfterAttackRoll{
		Source: source,
		Target: target,
		Result: result,
		Tags:   tags,
	}
}
