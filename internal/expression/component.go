package expression

import (
	"anvil/internal/tag"
)

var (
	Constant       = tag.FromString("Component.Type.Constant")
	DamageConstant = tag.FromString("Component.Type.Constant.Damage")
	Dice           = tag.FromString("Component.Type.Dice")
	D20            = tag.FromString("Component.Type.Dice.D20")
	DamageDice     = tag.FromString("Component.Type.Dice.Damage")
)

const (
	CriticalSuccess = 1
	CriticalFailure = -1
	CriticalNone    = 0
)

type Component struct {
	Type            tag.Tag
	Value           int
	Source          string
	Values          []int
	Times           int
	Sides           int
	HasAdvantage    []string
	HasDisadvantage []string
	Tags            tag.Container
	Components      []Component
	IsCritical      int
}

func (c *Component) hasRollModifier() bool {
	hasAdvantage := len(c.HasAdvantage) > 0
	hasDisadvantage := len(c.HasDisadvantage) > 0
	return hasAdvantage != hasDisadvantage
}

func (c Component) Clone() Component {
	var cloned []Component
	if len(c.Components) > 0 {
		cloned = make([]Component, len(c.Components))
		for i := range c.Components {
			cloned[i] = c.Components[i].Clone()
		}
	}

	return Component{
		Type:            c.Type,
		Source:          c.Source,
		Value:           c.Value,
		Times:           c.Times,
		Sides:           c.Sides,
		Values:          append([]int(nil), c.Values...),
		HasAdvantage:    append([]string(nil), c.HasAdvantage...),
		HasDisadvantage: append([]string(nil), c.HasDisadvantage...),
		Tags:            c.Tags,
		Components:      cloned,
		IsCritical:      c.IsCritical,
	}
}
