package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTagContainer_MatchTag(t *testing.T) {
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
			container := tag.ContainerFromString(tt.container)
			assert.Equal(t, tt.want, container.MatchTag(tag.FromString(tt.tag)))
		})
	}
}

func TestTagContainer_MatchAnyTag(t *testing.T) {
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
			container1 := tag.ContainerFromString(tt.container1)
			container2 := tag.ContainerFromString(tt.container2)
			assert.Equal(t, tt.want, container1.MatchAnyTag(container2))
		})
	}
}

func TestTagContainer_MatchAllTag(t *testing.T) {
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
			container1 := tag.ContainerFromString(tt.container1)
			var container2 tag.Container
			switch v := tt.container2.(type) {
			case string:
				container2 = tag.ContainerFromString(v)
			case []string:
				container2 = tag.ContainerFromStrings(v)
			}
			assert.Equal(t, tt.want, container1.MatchAllTag(container2))
		})
	}
}

func TestTagContainer_HasTag(t *testing.T) {
	tests := []struct {
		name      string
		container tag.Container
		tag       string
		want      bool
	}{
		{
			name:      "should match exact tags only - match",
			container: tag.ContainerFromStrings([]string{"ability.damage.fire"}),
			tag:       "ability.damage.fire",
			want:      true,
		},
		{
			name:      "should match exact tags only - no match for parent",
			container: tag.ContainerFromStrings([]string{"ability.damage.fire"}),
			tag:       "ability.damage",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.container.HasTag(tag.FromString(tt.tag)))
		})
	}
}

func TestTag_HasAnyTag(t *testing.T) {
	tests := []struct {
		name       string
		container1 tag.Container
		container2 tag.Container
		want       bool
	}{
		{
			name:       "should return true if any exact tag matches",
			container1: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: tag.ContainerFromStrings([]string{"Ability.Damage.ice", "ability.damage.poison"}),
			want:       true,
		},
		{
			name:       "should match exact - no match for parent",
			container1: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: tag.ContainerFromStrings([]string{"Ability.Damage"}),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.container1.HasAnyTag(tt.container2))
		})
	}
}

func TestTag_HasAllTags(t *testing.T) {
	tests := []struct {
		name       string
		container1 tag.Container
		container2 tag.Container
		want       bool
	}{
		{
			name:       "should return true if all exact tag matches",
			container1: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			want:       true,
		},
		{
			name:       "should return false if they don't match",
			container1: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.poison"}),
			want:       false,
		},
		{
			name:       "should match exact - no match for parent",
			container1: tag.ContainerFromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: tag.ContainerFromStrings([]string{"Ability.Damage"}),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.container1.HasAllTag(tt.container2))
		})
	}
}
