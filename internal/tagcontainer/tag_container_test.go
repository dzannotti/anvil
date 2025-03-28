package tagcontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTagContainer_Id(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialTag    string
		additionalTag string
		expectedId    string
	}{
		{
			name:          "should combine all tags returning id",
			initialTag:    "ability.damage.fire",
			additionalTag: "ability.damage.frost",
			expectedId:    "ability.damage.fire-ability.damage.frost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := FromString(tt.initialTag)
			container.AddTag(tag.FromString(tt.additionalTag))
			assert.Equal(t, tt.expectedId, container.ID())
		})
	}
}

func TestTagContainer_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialTag    string
		additionalTag string
		isEmpty       bool
	}{
		{
			name:          "should clone a tag container",
			initialTag:    "ability.damage.fire",
			additionalTag: "foo.bar",
			isEmpty:       false,
		},
		{
			name:          "should clone an empty container",
			isEmpty:       true,
			additionalTag: "foo.bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var container *TagContainer
			if tt.isEmpty {
				container = New()
			} else {
				container = FromString(tt.initialTag)
			}

			container2 := container.Clone()
			assert.Equal(t, container, container2)

			container2.AddTag(tag.FromString(tt.additionalTag))
			assert.NotEqual(t, container, container2)
		})
	}
}
