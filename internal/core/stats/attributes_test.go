package stats_test

import (
	"testing"

	"anvil/internal/core/stats"
	"anvil/internal/core/tags"

	"github.com/stretchr/testify/assert"
)

func TestAttribute_Value(t *testing.T) {
	attributes := stats.Attributes{
		Strength:     10,
		Dexterity:    11,
		Constitution: 12,
		Intelligence: 13,
		Wisdom:       14,
		Charisma:     15,
	}
	assert.Equal(t, 10, attributes.Value(tags.Strength))
	assert.Equal(t, 11, attributes.Value(tags.Dexterity))
	assert.Equal(t, 12, attributes.Value(tags.Constitution))
	assert.Equal(t, 13, attributes.Value(tags.Intelligence))
	assert.Equal(t, 14, attributes.Value(tags.Wisdom))
	assert.Equal(t, 15, attributes.Value(tags.Charisma))
}

func TestAttribute_Modifier(t *testing.T) {
	t.Run("should calculate positive modifiers correctly", func(t *testing.T) {
		assert.Equal(t, 2, stats.AttributeModifier(15))
		assert.Equal(t, 4, stats.AttributeModifier(18))
		assert.Equal(t, 5, stats.AttributeModifier(20))
	})

	t.Run("should calculate negative modifiers correctly", func(t *testing.T) {
		assert.Equal(t, -1, stats.AttributeModifier(8))
		assert.Equal(t, -2, stats.AttributeModifier(6))
		assert.Equal(t, -4, stats.AttributeModifier(3))
	})

	t.Run("should round down fractional modifiers", func(t *testing.T) {
		assert.Equal(t, 1, stats.AttributeModifier(13))
		assert.Equal(t, -2, stats.AttributeModifier(7))
	})
}
