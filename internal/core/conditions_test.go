package core

import (
	"testing"

	"anvil/internal/tag"
)

func TestConditions_Add(t *testing.T) {
	c := &Conditions{}
	testTag := tag.FromString("test")
	effect := &Effect{Name: "test"}

	c.Add(testTag, nil)
	if len(c.Sources[testTag]) != 0 {
		t.Error("Expected no effect to be added when source is nil")
	}

	c.Add(testTag, effect)
	if len(c.Sources[testTag]) != 1 {
		t.Error("Expected one effect to be added")
	}
	if c.Sources[testTag][0] != effect {
		t.Error("Expected added effect to match source")
	}

	effect2 := &Effect{Name: "test2"}
	c.Add(testTag, effect2)
	if len(c.Sources[testTag]) != 2 {
		t.Error("Expected two effects")
	}
}

func TestConditions_Remove(t *testing.T) {
	c := &Conditions{}
	testTag := tag.FromString("test")
	effect := &Effect{Name: "test"}

	removed := c.Remove(testTag, effect)
	if removed {
		t.Error("Expected false when removing from empty conditions")
	}

	// Add and remove specific effect
	c.Add(testTag, effect)
	removed = c.Remove(testTag, effect)
	if !removed {
		t.Error("Expected true when removing existing effect")
	}
	if len(c.Sources[testTag]) != 0 {
		t.Error("Expected effect to be removed")
	}

	// Test removing with nil source (should remove all)
	effect2 := &Effect{Name: "test2"}
	c.Add(testTag, effect)
	c.Add(testTag, effect2)
	removed = c.Remove(testTag, nil)
	if !removed {
		t.Error("Expected true when removing all effects")
	}
	if len(c.Sources[testTag]) != 0 {
		t.Error("Expected all effects to be removed")
	}
}

func TestConditions_Has(t *testing.T) {
	c := &Conditions{}
	testTag := tag.FromString("test")
	effect := &Effect{Name: "test"}

	// Test has on empty conditions
	if c.Has(testTag, effect) {
		t.Error("Expected Has to return false for empty conditions")
	}

	// Add effect and test Has
	c.Add(testTag, effect)
	if !c.Has(testTag, effect) {
		t.Error("Expected Has to return true for existing effect")
	}

	if !c.Has(testTag, nil) {
		t.Error("Expected Has to return true for nil source but present condition")
	}
	// Test Has with different effect
	effect2 := &Effect{Name: "test2"}
	if c.Has(testTag, effect2) {
		t.Error("Expected Has to return false for non-existing effect")
	}
}
