package state

type Type int

const (
	AttributeCalculationType Type = iota
	BeforeAttackRollType
)

type State interface {
	Type() Type
}
