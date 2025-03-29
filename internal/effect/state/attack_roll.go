package state

import (
	"anvil/internal/core/definition"
	"anvil/internal/expression"
	"anvil/internal/tagcontainer"
)

type BeforeAttackRoll struct {
	Source     definition.Creature
	Target     definition.Creature
	Expression *expression.Expression
	Tags       tagcontainer.TagContainer
}

func (s *BeforeAttackRoll) Type() Type {
	return BeforeAttackRollType
}

func NewBeforeAttackRoll(source definition.Creature, target definition.Creature, expr *expression.Expression, tags tagcontainer.TagContainer) *BeforeAttackRoll {
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
	Tags   tagcontainer.TagContainer
}

func (s *AfterAttackRoll) Type() Type {
	return AfterAttackRollType
}

func NewAfterAttackRoll(source definition.Creature, target definition.Creature, result *expression.Expression, tags tagcontainer.TagContainer) *AfterAttackRoll {
	return &AfterAttackRoll{
		Source: source,
		Target: target,
		Result: result,
		Tags:   tags,
	}
}
