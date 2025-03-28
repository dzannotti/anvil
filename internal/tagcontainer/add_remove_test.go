package tagcontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTagContainer_AddTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialTagStr  string
		tagToAdd       string
		expectedLength int
	}{
		{
			name:           "can add tag",
			initialTagStr:  "",
			tagToAdd:       "ability.damage.fire",
			expectedLength: 1,
		},
		{
			name:           "does not add duplicates",
			initialTagStr:  "ability.damage.fire",
			tagToAdd:       "ability.damage.fire",
			expectedLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var container *TagContainer
			if tt.initialTagStr == "" {
				container = New()
			} else {
				container = FromString(tt.initialTagStr)
			}
			container.AddTag(tag.FromString(tt.tagToAdd))
			assert.Equal(t, tt.expectedLength, len(container.tags))
			assert.True(t, tag.FromString(tt.tagToAdd).MatchExact(container.tags[0]))
		})
	}
}

func TestTagContainer_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialTagStr  string
		tagsToAdd      []string
		expectedLength int
		expectedTags   []string
	}{
		{
			name:           "can add tag containers",
			initialTagStr:  "ability.damage.fire",
			tagsToAdd:      []string{"ability.damage.ice"},
			expectedLength: 2,
			expectedTags:   []string{"ability.damage.fire", "ability.damage.ice"},
		},
		{
			name:           "does not add duplicates",
			initialTagStr:  "ability.damage.fire",
			tagsToAdd:      []string{"ability.damage.fire", "ability.damage.ice"},
			expectedLength: 2,
			expectedTags:   []string{"ability.damage.fire", "ability.damage.ice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container1 := FromString(tt.initialTagStr)
			container2 := FromStrings(tt.tagsToAdd)
			container1.Add(*container2)
			assert.Equal(t, tt.expectedLength, len(container1.tags))
			for _, expectedTag := range tt.expectedTags {
				assert.True(t, container1.HasTag(tag.FromString(expectedTag)))
			}
		})
	}
}

func TestTagContainer_RemoveTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialTagStr  string
		tagToRemove    string
		expectedLength int
	}{
		{
			name:           "should remove tag",
			initialTagStr:  "ability.damage.fire",
			tagToRemove:    "ability.damage.fire",
			expectedLength: 0,
		},
		{
			name:           "should not remove relative tags",
			initialTagStr:  "ability.damage.fire",
			tagToRemove:    "ability.damage",
			expectedLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := New()
			container.AddTag(tag.FromString(tt.initialTagStr))
			container.RemoveTag(tag.FromString(tt.tagToRemove))
			assert.Equal(t, tt.expectedLength, len(container.tags))
		})
	}
}
