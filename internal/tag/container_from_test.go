package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestContainer_From(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func() tag.Container
		expected []tag.Tag
	}{
		{
			name: "creates valid empty container",
			setup: func() tag.Container {
				return tag.Container{}
			},
			expected: []tag.Tag{},
		},
		{
			name: "creates valid container from string",
			setup: func() tag.Container {
				return tag.ContainerFromString("ability.damage.fire")
			},
			expected: []tag.Tag{tag.FromString("ability.damage.fire")},
		},
		{
			name: "creates valid container from slice of string",
			setup: func() tag.Container {
				return tag.ContainerFromStrings([]string{"ability.damage.fire", "ability.damage.ice"})
			},
			expected: []tag.Tag{
				tag.FromString("ability.damage.fire"),
				tag.FromString("ability.damage.ice"),
			},
		},
		{
			name: "creates valid container from tag",
			setup: func() tag.Container {
				return tag.ContainerFromTag(tag.FromString("ability.damage.fire"))
			},
			expected: []tag.Tag{tag.FromString("ability.damage.fire")},
		},
		{
			name: "creates valid container from slice of tag",
			setup: func() tag.Container {
				return tag.ContainerFromTag(
					tag.FromString("ability.damage.fire"),
					tag.FromString("ability.damage.ice"),
				)
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
			assert.Equal(t, len(tt.expected), len(container.Strings()))
			for i, expectedTag := range tt.expected {
				assert.True(t, expectedTag.MatchExact(tag.FromString(container.Strings()[i])))
			}
		})
	}
}
