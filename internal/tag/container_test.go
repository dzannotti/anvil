package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestContainer_Constructor(t *testing.T) {
	t.Run("creates from string", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})

	t.Run("creates from tags", func(t *testing.T) {
		c := tag.ContainerFromTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})
}

func TestContainer_HasTag(t *testing.T) {
	t.Run("returns true if we have tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.True(t, c.HasTag(tag.FromString("ability.damage")))
	})

	t.Run("returns false if we do not have tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
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
		c := tag.ContainerFromString("ability.damage")
		c.AddTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{"ability.damage"}, c.AsStrings())
	})
}

func TestContainer_RemoveTag(t *testing.T) {
	t.Run("should remove existing tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		c.RemoveTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{}, c.AsStrings())
	})

	t.Run("should not panic when removing non existing tag", func(t *testing.T) {
		c := tag.Container{}
		c.RemoveTag(tag.FromString("ability.damage"))
		assert.Equal(t, []string{}, c.AsStrings())
	})

	t.Run("should remove multiple tags at once", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage", "item.weapon", "spell.fire")
		c.RemoveTag(tag.FromString("ability.damage"), tag.FromString("spell.fire"))
		assert.Equal(t, []string{"item.weapon"}, c.AsStrings())
	})
}

func TestContainer_HasAny(t *testing.T) {
	t.Run("should match any container tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.True(t, c.HasAny(tag.ContainerFromString("item.weapon", "ability.damage")))
	})

	t.Run("should not match missing tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.False(t, c.HasAny(tag.ContainerFromString("item.weapon")))
	})
}

func TestContainer_HasAll(t *testing.T) {
	t.Run("should match all container tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage", "item.weapon")
		assert.True(t, c.HasAll(tag.ContainerFromString("item.weapon", "ability.damage")))
	})

	t.Run("should not match missing tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.False(t, c.HasAll(tag.ContainerFromString("item.weapon")))
	})

	t.Run("should not match partial container tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.False(t, c.HasAll(tag.ContainerFromString("item.weapon", "ability.damage")))
	})
}

func TestContainer_MatchTag(t *testing.T) {
	t.Run("matches more general tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage.fire")
		assert.True(t, c.MatchTag(tag.FromString("ability.damage")))
	})

	t.Run("does not matches more specific tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.False(t, c.MatchTag(tag.FromString("ability.damage.fire")))
	})

	t.Run("does not match partial component", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage.fire")
		assert.False(t, c.MatchTag(tag.FromString("ability.dam")))
	})
}

func TestContainer_MatchAnyTags(t *testing.T) {
	t.Run("matches more general tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage.fire", "item.weapon")
		assert.True(t, c.MatchAny(tag.ContainerFromString("ability.melee", "item.weapon")))
	})

	t.Run("does not matches more specific tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.spell", "ability.damage")
		assert.False(t, c.MatchAny(tag.ContainerFromString("ability.damage.fire")))
	})
}

func TestContainer_MatchAllTag(t *testing.T) {
	t.Run("requires all tags matching", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage.fire", "item.weapon")
		assert.True(t, c.MatchAll(tag.ContainerFromString("ability.damage.fire", "item.weapon")))
	})

	t.Run("does not match if one differs", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage.fire", "item.weapon")
		assert.False(t, c.MatchAll(tag.ContainerFromString("ability.damage.fire", "item.armor")))
	})
}

func TestContainer_IsEmpty(t *testing.T) {
	t.Run("returns true if empty", func(t *testing.T) {
		c := tag.Container{}
		assert.True(t, c.IsEmpty())
	})

	t.Run("returns false if not empty", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage")
		assert.False(t, c.IsEmpty())
	})
}

func TestContainer_ID(t *testing.T) {
	t.Run("returns tag id from tag", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage", "item.weapon")
		assert.Equal(t, "ability.damage-item.weapon", c.ID())
	})

	t.Run("the order of tags is idempotent", func(t *testing.T) {
		c1 := tag.ContainerFromString("ability.damage", "item.weapon")
		c2 := tag.ContainerFromString("item.weapon", "ability.damage")
		assert.Equal(t, c1.ID(), c2.ID())
	})
}

func TestContainer_Len(t *testing.T) {
	t.Run("returns zero for empty container", func(t *testing.T) {
		c := tag.Container{}
		assert.Equal(t, 0, c.Len())
	})

	t.Run("returns correct count", func(t *testing.T) {
		c := tag.ContainerFromString("ability.damage", "item.weapon")
		assert.Equal(t, 2, c.Len())
	})
}

func TestContainer_Clone(t *testing.T) {
	t.Run("creates independent copy", func(t *testing.T) {
		original := tag.ContainerFromString("ability.damage")
		clone := original.Clone()

		clone.AddTag(tag.FromString("item.weapon"))

		assert.Equal(t, 1, original.Len())
		assert.Equal(t, 2, clone.Len())
	})
}

func TestContainer_Add(t *testing.T) {
	t.Run("merges tags from other containers", func(t *testing.T) {
		c1 := tag.ContainerFromString("ability.damage")
		c2 := tag.ContainerFromString("item.weapon")
		c3 := tag.ContainerFromString("spell.fire")

		c1.Add(c2, c3)

		assert.Equal(t, 3, c1.Len())
		assert.True(t, c1.HasTag(tag.FromString("ability.damage")))
		assert.True(t, c1.HasTag(tag.FromString("item.weapon")))
		assert.True(t, c1.HasTag(tag.FromString("spell.fire")))
	})

	t.Run("prevents duplicates when merging", func(t *testing.T) {
		c1 := tag.ContainerFromString("ability.damage")
		c2 := tag.ContainerFromString("ability.damage", "item.weapon")

		c1.Add(c2)

		assert.Equal(t, 2, c1.Len())
	})
}
