package creature

import (
	"anvil/internal/core/tags"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttribute_Value(t *testing.T) {
	attributes := NewAttributes(AttributeValues{Strength: 10, Dexterity: 11, Constitution: 12, Intelligence: 13, Wisdom: 14, Charisma: 15})
	assert.Equal(t, 10, attributes.Value(tags.Strength))
	assert.Equal(t, 11, attributes.Value(tags.Dexterity))
	assert.Equal(t, 12, attributes.Value(tags.Constitution))
	assert.Equal(t, 13, attributes.Value(tags.Intelligence))
	assert.Equal(t, 14, attributes.Value(tags.Wisdom))
	assert.Equal(t, 15, attributes.Value(tags.Charisma))
}

func TestAttribute_Modifier(t *testing.T) {
	t.Parallel()
	t.Run("should calculate positive modifiers correctly", func(t *testing.T) {
		assert.Equal(t, 2, AttributeModifier(15))
		assert.Equal(t, 4, AttributeModifier(18))
		assert.Equal(t, 5, AttributeModifier(20))
	})

	t.Run("should calculate negative modifiers correctly", func(t *testing.T) {
		assert.Equal(t, -1, AttributeModifier(8))
		assert.Equal(t, -2, AttributeModifier(6))
		assert.Equal(t, -4, AttributeModifier(3))
	})

	t.Run("should round down fractional modifiers", func(t *testing.T) {
		assert.Equal(t, 1, AttributeModifier(13))
		assert.Equal(t, -2, AttributeModifier(7))
	})
}
