package expression

import "anvil/internal/tag"

type ComponentKind string

const (
	ComponentKindConstant ComponentKind = "constant"
	ComponentKindDice     ComponentKind = "dice"
	ComponentKindD20      ComponentKind = "d20"
)

type Component interface {
	Kind() ComponentKind
	Tags() tag.Container
	Value() int
	Source() string
	Components() []Component
	Clone() Component
	Evaluate(ctx *Context) int
	Expected() int
}
