package creature

import (
	"anvil/internal/core/tags"
	"anvil/internal/tag"
	"math"
)

type Attributes struct {
	values map[string]int
}

type AttributeValues struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

func NewAttributes(values AttributeValues) Attributes {
	return Attributes{
		values: map[string]int{
			tags.Strength.String():     values.Strength,
			tags.Dexterity.String():    values.Dexterity,
			tags.Constitution.String(): values.Constitution,
			tags.Intelligence.String(): values.Intelligence,
			tags.Wisdom.String():       values.Wisdom,
			tags.Charisma.String():     values.Charisma,
		},
	}
}

func (a Attributes) Value(tag tag.Tag) int {
	val, exists := a.values[tag.String()]
	if exists {
		return val
	}
	return 0
}

func AttributeModifier(value int) int {
	if value < 1 {
		return -5 // minimum modifier for a score of 1
	}
	return int(math.Floor(float64(value-10) / 2.0))
}
