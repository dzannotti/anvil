package tagcontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTagContainer_MatchTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		container string
		tag       string
		want      bool
	}{
		{
			name:      "should match exact tag",
			container: "ability.damage.fire",
			tag:       "ability.damage.fire",
			want:      true,
		},
		{
			name:      "should match parent tag",
			container: "ability.damage.fire",
			tag:       "ability.damage",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := FromString(tt.container)
			assert.Equal(t, tt.want, container.MatchTag(tag.FromString(tt.tag)))
		})
	}
}

func TestTagContainer_MatchAnyTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		container1 string
		container2 string
		want       bool
	}{
		{
			name:       "should return true if any tag matches exactly",
			container1: "ability.damage.fire",
			container2: "ability.damage.fire",
			want:       true,
		},
		{
			name:       "should return true if any parent tag matches",
			container1: "ability.damage.fire",
			container2: "ability.damage",
			want:       true,
		},
		{
			name:       "should return false if only child tags exist",
			container1: "ability.damage",
			container2: "ability.damage.fire",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container1 := FromString(tt.container1)
			container2 := FromString(tt.container2)
			assert.Equal(t, tt.want, container1.MatchAnyTag(*container2))
		})
	}
}

func TestTagContainer_MatchAllTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		container1 string
		container2 interface{}
		want       bool
	}{
		{
			name:       "should return true if all tags match exactly",
			container1: "ability.damage.fire",
			container2: "ability.damage.fire",
			want:       true,
		},
		{
			name:       "should return true if all parent tags match",
			container1: "ability.damage.fire",
			container2: "ability.damage",
			want:       true,
		},
		{
			name:       "should return false if any tag is missing",
			container1: "ability.damage.fire",
			container2: []string{"ability.damage", "status.burning"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container1 := FromString(tt.container1)
			var container2 *TagContainer
			switch v := tt.container2.(type) {
			case string:
				container2 = FromString(v)
			case []string:
				container2 = FromStrings(v)
			}
			assert.Equal(t, tt.want, container1.MatchAllTag(*container2))
		})
	}
}
