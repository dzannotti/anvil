package expression

import "anvil/internal/tag"

type ComponentType string

const (
	TypeConstant      ComponentType = "constant"
	TypeDamageConstant ComponentType = "constant-damage"
	TypeDice          ComponentType = "dice"
	TypeDice20        ComponentType = "dice-20"
	TypeDamageDice    ComponentType = "dice-damage"
)

type Component struct {
	Type            ComponentType
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

func (c *Component) shouldModifyRoll() bool {
	return (len(c.HasAdvantage) > 0) != (len(c.HasDisadvantage) > 0)
}

func (c Component) Clone() Component {
	var cloned []Component
	if c.Components != nil {
		cloned = make([]Component, len(c.Components))
		for i := range c.Components {
			cloned[i] = c.Components[i].Clone()
		}
	}

	var values []int
	if c.Values != nil {
		values = append(make([]int, 0), c.Values...)
	}

	var hasAdvantage []string
	if c.HasAdvantage != nil {
		hasAdvantage = append(make([]string, 0), c.HasAdvantage...)
	}

	var hasDisadvantage []string
	if c.HasDisadvantage != nil {
		hasDisadvantage = append(make([]string, 0), c.HasDisadvantage...)
	}

	return Component{
		Type:            c.Type,
		Source:          c.Source,
		Value:           c.Value,
		Times:           c.Times,
		Sides:           c.Sides,
		Values:          values,
		HasAdvantage:    hasAdvantage,
		HasDisadvantage: hasDisadvantage,
		Tags:            c.Tags,
		Components:      cloned,
		IsCritical:      c.IsCritical,
	}
}

