package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestConditions_Add(t *testing.T) {
	c := &Conditions{}
	testTag := tag.FromString("test")
	effect := &Effect{Name: "test"}

	assert.Panics(t, func() {
		c.Add(testTag, nil)
	}, "Expected conditions to panic when source is nil")

	c.Add(testTag, effect)
	assert.Len(t, c.Sources[testTag], 1, "Expected one effect to be added")
	assert.Equal(t, effect, c.Sources[testTag][0], "Expected added effect to match source")

	effect2 := &Effect{Name: "test2"}
	c.Add(testTag, effect2)
	assert.Len(t, c.Sources[testTag], 2, "Expected two effects")
}

func TestConditions_Remove(t *testing.T) {
	c := &Conditions{}
	testTag := tag.FromString("test")
	effect := &Effect{Name: "test"}

	removed := c.Remove(testTag, effect)
	assert.False(t, removed, "Expected false when removing from empty conditions")

	// Add and remove specific effect
	c.Add(testTag, effect)
	removed = c.Remove(testTag, effect)
	assert.True(t, removed, "Expected true when removing existing effect")
	assert.Len(t, c.Sources[testTag], 0, "Expected effect to be removed")

	// Test removing with nil source (should remove all)
	effect2 := &Effect{Name: "test2"}
	c.Add(testTag, effect)
	c.Add(testTag, effect2)
	removed = c.Remove(testTag, nil)
	assert.True(t, removed, "Expected true when removing all effects")
	assert.Len(t, c.Sources[testTag], 0, "Expected all effects to be removed")
}

func TestConditions_Has(t *testing.T) {
	c := &Conditions{}
	testTag := tag.FromString("test")
	effect := &Effect{Name: "test"}

	// Test has on empty conditions
	assert.False(t, c.Has(testTag, effect), "Expected Has to return false for empty conditions")

	// Add effect and test Has
	c.Add(testTag, effect)
	assert.True(t, c.Has(testTag, effect), "Expected Has to return true for existing effect")

	assert.True(t, c.Has(testTag, nil), "Expected Has to return true for nil source but present condition")
	// Test Has with different effect
	effect2 := &Effect{Name: "test2"}
	assert.False(t, c.Has(testTag, effect2), "Expected Has to return false for non-existing effect")
}
