package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestContainer_Constructor(t *testing.T) {
	t.Run("creates from string", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})

	t.Run("creates from tags", func(t *testing.T) {
		c := tag.NewContainer(tag.FromString("ability.damage"))
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})

	t.Run("creates from container", func(t *testing.T) {
		c1 := tag.NewContainer(tag.FromString("ability.damage"))
		c2 := tag.NewContainerFromContainer(c1)
		assert.Equal(t, c1, c2)
	})
}

func TestContainer_HasTag(t *testing.T) {
	t.Run("returns true if we have tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.True(t, c.HasTag(tag.FromString("ability.damage")))
	})

	t.Run("returns false if we do not have tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.False(t, c.HasTag(tag.FromString("item.weapon")))
	})
}

func TestContainer_AddTag(t *testing.T) {
	t.Run("can add tag", func(t *testing.T) {
		c := tag.Container{}
		c.AddTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})

	t.Run("does not add existing tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		c.AddTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})
}

func TestContainer_RemoveTag(t *testing.T) {
	t.Run("should remove existing tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		c.RemoveTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{}, c.AsStrings())
	})

	t.Run("should not panic when removing non existing tag", func(t *testing.T) {
		c := tag.Container{}
		c.RemoveTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{}, c.AsStrings())
	})
}

func TestContainer_HasAnyTag(t *testing.T) {
	t.Run("should match any container tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.True(t, c.HasAny(tag.NewContainerFromString("item.weapon", "ability.damage")))
	})

	t.Run("should not match missing tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.False(t, c.HasAny(tag.NewContainerFromString("item.weapon")))
	})
}

func TestContainer_HasAllTag(t *testing.T) {
	t.Run("should match all container tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage", "item.weapon")
		assert.True(t, c.HasAll(tag.NewContainerFromString("item.weapon", "ability.damage")))
	})

	t.Run("should not match missing tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.False(t, c.HasAll(tag.NewContainerFromString("item.weapon")))
	})

	t.Run("should not match partial container tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.False(t, c.HasAll(tag.NewContainerFromString("item.weapon", "ability.damage")))
	})
}

func TestContainer_MatchTag(t *testing.T) {
	t.Run("matches more general tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage.fire")
		assert.True(t, c.MatchTag(tag.FromString("ability.damage")))
	})

	t.Run("does not matches more specific tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.False(t, c.MatchTag(tag.FromString("ability.damage.fire")))
	})

	t.Run("does not match partial component", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage.fire")
		assert.False(t, c.MatchTag(tag.FromString("ability.dam")))
	})
}

func TestContainer_MatchAnyTags(t *testing.T) {
	t.Run("matches more general tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage.fire", "item.weapon")
		assert.True(t, c.MatchAny(tag.NewContainerFromString("ability.melee", "item.weapon")))
	})

	t.Run("does not matches more specific tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.spell", "ability.damage")
		assert.False(t, c.MatchAny(tag.NewContainerFromString("ability.damage.fire")))
	})
}

func TestContainer_MatchAllTag(t *testing.T) {
	t.Run("requires all tags matching", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage.fire", "item.weapon")
		assert.True(t, c.MatchAll(tag.NewContainerFromString("ability.damage.fire", "item.weapon")))
	})

	t.Run("does not match if one differs", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage.fire", "item.weapon")
		assert.False(t, c.MatchAll(tag.NewContainerFromString("ability.damage.fire", "item.armor")))
	})
}

func TestContainer_IsEmpty(t *testing.T) {
	t.Run("returns true if empty", func(t *testing.T) {
		c := tag.Container{}
		assert.True(t, c.IsEmpty())
	})

	t.Run("returns false if not empty", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage")
		assert.False(t, c.IsEmpty())
	})
}

func TestContainer_ID(t *testing.T) {
	t.Run("returns tag id from tag", func(t *testing.T) {
		c := tag.NewContainerFromString("ability.damage", "item.weapon")
		assert.Equal(t, "ability.damage-item.weapon", c.ID())
	})

	t.Run("the order of tags is idempotent", func(t *testing.T) {
		c1 := tag.NewContainerFromString("ability.damage", "item.weapon")
		c2 := tag.NewContainerFromString("item.weapon", "ability.damage")
		assert.Equal(t, c1.ID(), c2.ID())
	})
}
