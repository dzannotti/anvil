package stats

import (
	"math"

	"anvil/internal/core/tags"
	"anvil/internal/tag"
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
	case tags.AttributeStrength:
		return a.Strength
	case tags.AttributeDexterity:
		return a.Dexterity
	case tags.AttributeConstitution:
		return a.Constitution
	case tags.AttributeIntelligence:
		return a.Intelligence
	case tags.AttributeWisdom:
		return a.Wisdom
	case tags.AttributeCharisma:
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
