package tag_test

import (
	"anvil/internal/tag"
	"testing"

	"github.com/stretchr/testify/assert"
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
			container := tag.ContainerFromString(tt.initialTagStr)
			container.AddTag(tag.FromString(tt.tagToAdd))
			assert.Equal(t, tt.expectedLength, len(container.Strings()))
			assert.True(t, tag.FromString(tt.tagToAdd).MatchExact(tag.FromString(container.Strings()[0])))
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
			container1 := tag.ContainerFromString(tt.initialTagStr)
			container2 := tag.ContainerFromStrings(tt.tagsToAdd)
			container1.Add(container2)
			assert.Equal(t, tt.expectedLength, len(container1.Strings()))
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
			container := tag.NewContainer()
			container.AddTag(tag.FromString(tt.initialTagStr))
			container.RemoveTag(tag.FromString(tt.tagToRemove))
			assert.Equal(t, tt.expectedLength, len(container.Strings()))
		})
	}
}
