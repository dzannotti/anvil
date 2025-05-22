package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTag_Constructor(t *testing.T) {
	t.Run("creates from string", func(t *testing.T) {
		tt := tag.FromString("ability.damage")
		assert.Equal(t, "ability.damage", tt.AsString())
	})

	t.Run("removes leading and trailing space", func(t *testing.T) {
		tt := tag.FromString("  ability.  damage  ")
		assert.Equal(t, "ability.damage", tt.AsString())
	})

	t.Run("removes any special character", func(t *testing.T) {
		tt := tag.FromString("ability.dama@$%&ge")
		assert.Equal(t, "ability.damage", tt.AsString())
	})

	t.Run("removes any extra dot separation", func(t *testing.T) {
		tt := tag.FromString(".ability..damage.")
		assert.Equal(t, "ability.damage", tt.AsString())
	})

	t.Run("removes any unicode character", func(t *testing.T) {
		tt := tag.FromString("ability.damage.🔥")
		assert.Equal(t, "ability.damage", tt.AsString())
	})
}

func TestTag_MatchExact(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		tt := tag.FromString("ability.damage")
		assert.True(t, tt.MatchExact(tag.FromString("ability.damage")))
	})
	
	t.Run("not-exact match", func(t *testing.T) {
		tt := tag.FromString("ability.damage")
		assert.False(t, tt.MatchExact(tag.FromString("ability")))
		assert.False(t, tt.MatchExact(tag.FromString("abiity.damage.fire")))
	})
}

func TestTag_Match(t *testing.T) {
	t.Run("matches more general tag", func(t *testing.T) {
		tt := tag.FromString("ability.damage.fire")
		assert.True(t, tt.Match(tag.FromString("ability.damage")))
	})

	t.Run("does not matches more specific tag", func(t *testing.T) {
		tt := tag.FromString("ability.damage")
		assert.False(t, tt.Match(tag.FromString("ability.damage.fire")))
	})

	t.Run("does not match partial component", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage.fire")
		assert.False(t, c.MatchTag(tag.FromString("ability.dam")))
	})
}

func TestTag_AsStrings(t *testing.T) {
	t.Run("splits tag into strings", func(t *testing.T) {
		tt := tag.FromString("ability.damage.fire")
		assert.Equal(t, []string{"ability", "damage", "fire"}, tt.AsStrings())
	})
}
