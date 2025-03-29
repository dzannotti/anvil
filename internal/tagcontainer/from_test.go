package tagcontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTagContainer_From(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func() TagContainer
		expected []tag.Tag
	}{
		{
			name: "creates valid empty container",
			setup: func() TagContainer {
				return New()
			},
			expected: []tag.Tag{},
		},
		{
			name: "creates valid container from string",
			setup: func() TagContainer {
				return FromString("ability.damage.fire")
			},
			expected: []tag.Tag{tag.FromString("ability.damage.fire")},
		},
		{
			name: "creates valid container from slice of string",
			setup: func() TagContainer {
				return FromStrings([]string{"ability.damage.fire", "ability.damage.ice"})
			},
			expected: []tag.Tag{
				tag.FromString("ability.damage.fire"),
				tag.FromString("ability.damage.ice"),
			},
		},
		{
			name: "creates valid container from tag",
			setup: func() TagContainer {
				return FromTag(tag.FromString("ability.damage.fire"))
			},
			expected: []tag.Tag{tag.FromString("ability.damage.fire")},
		},
		{
			name: "creates valid container from slice of tag",
			setup: func() TagContainer {
				return FromTags([]tag.Tag{
					tag.FromString("ability.damage.fire"),
					tag.FromString("ability.damage.ice"),
				})
			},
			expected: []tag.Tag{
				tag.FromString("ability.damage.fire"),
				tag.FromString("ability.damage.ice"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := tt.setup()
			assert.Equal(t, len(tt.expected), len(container.tags))
			for i, expectedTag := range tt.expected {
				assert.True(t, expectedTag.MatchExact(container.tags[i]))
			}
		})
	}
}
