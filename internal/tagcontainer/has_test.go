package tagcontainer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/tag"
)

func TestTagContainer_HasTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		container *TagContainer
		tag       string
		want      bool
	}{
		{
			name:      "should match exact tags only - match",
			container: FromStrings([]string{"ability.damage.fire"}),
			tag:       "ability.damage.fire",
			want:      true,
		},
		{
			name:      "should match exact tags only - no match for parent",
			container: FromStrings([]string{"ability.damage.fire"}),
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
	t.Parallel()

	tests := []struct {
		name       string
		container1 *TagContainer
		container2 *TagContainer
		want       bool
	}{
		{
			name:       "should return true if any exact tag matches",
			container1: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: FromStrings([]string{"Ability.Damage.ice", "ability.damage.poison"}),
			want:       true,
		},
		{
			name:       "should match exact - no match for parent",
			container1: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: FromStrings([]string{"Ability.Damage"}),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.container1.HasAnyTag(*tt.container2))
		})
	}
}

func TestTag_HasAllTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		container1 *TagContainer
		container2 *TagContainer
		want       bool
	}{
		{
			name:       "should return true if all exact tag matches",
			container1: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			want:       true,
		},
		{
			name:       "should return false if they don't match",
			container1: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.poison"}),
			want:       false,
		},
		{
			name:       "should match exact - no match for parent",
			container1: FromStrings([]string{"Ability.Damage.Fire", "ability.damage.ice"}),
			container2: FromStrings([]string{"Ability.Damage"}),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.container1.HasAllTag(*tt.container2))
		})
	}
}
