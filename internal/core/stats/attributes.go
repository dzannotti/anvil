package stats

import (
	"anvil/internal/core/tags"
	"anvil/internal/tag"
	"math"
)

type Attributes struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

func (a Attributes) Value(tag tag.Tag) int {
	switch tag {
	case tags.Strength:
		return a.Strength
	case tags.Dexterity:
		return a.Dexterity
	case tags.Constitution:
		return a.Constitution
	case tags.Intelligence:
		return a.Intelligence
	case tags.Wisdom:
		return a.Wisdom
	case tags.Charisma:
		return a.Charisma
	}
	return 0
}

func AttributeModifier(value int) int {
	if value < 1 {
		return -5 // minimum modifier for a score of 1
	}
	return int(math.Floor(float64(value-10) / 2.0))
}
