package state

import (
	"anvil/internal/expression"
	"anvil/internal/tagcontainer"
)

type BeforeAttackRoll struct {
	Expression expression.Expression
	Tags       tagcontainer.TagContainer
}

func (s *BeforeAttackRoll) Type() Type {
	return BeforeAttackRollType
}
