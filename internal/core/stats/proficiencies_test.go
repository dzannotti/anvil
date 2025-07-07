package stats_test

import (
	"testing"

	"anvil/internal/core/stats"
	"anvil/internal/tag"

	"github.com/stretchr/testify/assert"
)

func TestProficiencies_NewProficiencies(t *testing.T) {
	prof := stats.Proficiencies{Bonus: 10}
	assert.False(t, prof.Has(tag.NewContainerFromString("any")), "New proficiencies should have no skills")
}

func TestProficiencies_Add(t *testing.T) {
	t.Run("single tag", func(t *testing.T) {
		prof := stats.Proficiencies{Bonus: 10}
		prof.Add(tag.FromString("test"))
		assert.True(t, prof.Has(tag.NewContainerFromString("test")))
	})
	t.Run("multiple tags", func(t *testing.T) {
		prof := stats.Proficiencies{Bonus: 10}
		prof.Add(tag.FromString("test1"))
		prof.Add(tag.FromString("test2"))
		assert.True(t, prof.Has(tag.NewContainerFromString("test1")))
		assert.True(t, prof.Has(tag.NewContainerFromString("test2")))
	})
}

func TestProficiencies_Has(t *testing.T) {
	t.Run("should return true for added proficiency", func(t *testing.T) {
		prof := stats.Proficiencies{Bonus: 2}
		prof.Add(tag.FromString("test"))
		assert.True(t, prof.Has(tag.NewContainerFromString("test")))
	})

	t.Run("should return false for non-added proficiency", func(t *testing.T) {
		prof := stats.Proficiencies{Bonus: 2}
		assert.False(t, prof.Has(tag.NewContainerFromString("test")))
	})

	t.Run("should return true for hierarchical tag match (child->parent)", func(t *testing.T) {
		prof := stats.Proficiencies{Bonus: 2}
		prof.Add(tag.FromString("weapon.martial"))
		assert.True(t, prof.Has(tag.NewContainerFromString("weapon.martial.sword")))
	})

	t.Run("should return false for hierarchical tag match (parent->child)", func(t *testing.T) {
		prof := stats.Proficiencies{Bonus: 2}
		prof.Add(tag.FromString("weapon.martial.sword"))
		assert.False(t, prof.Has(tag.NewContainerFromString("weapon.martial")))
	})
}
