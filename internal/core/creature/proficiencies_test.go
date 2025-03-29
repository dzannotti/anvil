package creature

import (
	"anvil/internal/tag"
	"anvil/internal/tagcontainer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProficiencies_NewProficiencies(t *testing.T) {
	prof := NewProficiencies(10)
	assert.Equal(t, 10, prof.bonus)
	assert.False(t, prof.Has(tagcontainer.FromString("any")), "New proficiencies should have no skills")
}

func TestProficiencies_Add(t *testing.T) {
	t.Run("single tag", func(t *testing.T) {
		prof := NewProficiencies(10)
		prof.Add(tag.FromString("test"))
		assert.True(t, prof.Has(tagcontainer.FromString("test")))
	})
	t.Run("multiple tags", func(t *testing.T) {
		prof := NewProficiencies(10)
		prof.Add(tag.FromString("test1"))
		prof.Add(tag.FromString("test2"))
		assert.True(t, prof.Has(tagcontainer.FromString("test1")))
		assert.True(t, prof.Has(tagcontainer.FromString("test2")))
	})
}

func TestProficiencies_Has(t *testing.T) {
	t.Parallel()
	t.Run("should return true for added proficiency", func(t *testing.T) {
		proficiencies := NewProficiencies(2)
		proficiencies.Add(tag.FromString("test"))
		assert.True(t, proficiencies.Has(tagcontainer.FromString("test")))
	})

	t.Run("should return false for non-added proficiency", func(t *testing.T) {
		proficiencies := NewProficiencies(2)
		assert.False(t, proficiencies.Has(tagcontainer.FromString("test")))
	})

	t.Run("should return true for hierarchical tag match (child->parent)", func(t *testing.T) {
		proficiencies := NewProficiencies(2)
		proficiencies.Add(tag.FromString("weapon.martial"))
		assert.True(t, proficiencies.Has(tagcontainer.FromString("weapon.martial.sword")))
	})

	t.Run("should return false for hierarchical tag match (parent->child)", func(t *testing.T) {
		proficiencies := NewProficiencies(2)
		proficiencies.Add(tag.FromString("weapon.martial.sword"))
		assert.False(t, proficiencies.Has(tagcontainer.FromString("weapon.martial")))
	})
}