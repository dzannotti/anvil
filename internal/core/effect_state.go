package core

import (
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type PreAttackRoll struct {
	Source     *Actor
	Target     *Actor
	Expression *expression.Expression
	Tags       tag.Container
}

type PostAttackRoll struct {
	Source *Actor
	Target *Actor
	Result *expression.Expression
	Tags   tag.Container
}

type AttributeCalculation struct {
	Source     *Actor
	Expression *expression.Expression
	Attribute  tag.Tag
}

type PreTakeDamage struct {
	Expression *expression.Expression
	Source     *Actor
}

type PostTakeDamage struct {
	Result       *expression.Expression
	Source       *Actor
	ActualDamage int
}

type PreDamageRoll struct {
	Expression *expression.Expression
	Source     *Actor
	Tags       tag.Container
	Critical   bool
}

type PostDamageRoll struct {
	Result   *expression.Expression
	Source   *Actor
	Tags     tag.Container
	Critical bool
}

type PreSavingThrow struct {
	Expression      *expression.Expression
	Source          *Actor
	Attribute       tag.Tag
	DifficultyClass int
}

type PostSavingThrow struct {
	Result          *expression.Expression
	Source          *Actor
	Attribute       tag.Tag
	DifficultyClass int
}

type AttributeChanged struct {
	Source    *Actor
	Attribute tag.Tag
	OldValue  int
	Value     int
}

type ConditionChanged struct {
	Source    *Actor
	Condition tag.Tag
	From      *Effect
}

type TurnStarted struct {
	Source *Actor
}

type TurnEnded struct {
	Source *Actor
}

type PreMoveStep struct {
	Source  *Actor
	Action  Action
	From    grid.Position
	To      grid.Position
	CanMove bool
}
