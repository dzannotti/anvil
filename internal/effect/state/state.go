package state

type Type int

const (
	AttributeCalculationType Type = iota
	BeforeAttackRollType
	AfterAttackRollType
)

type State interface {
	Type() Type
}
