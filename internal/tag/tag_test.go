package tag_test

import (
	"anvil/internal/tag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "creates valid tag with standard input",
			input:    "ability.damage.fire",
			expected: "ability.damage.fire",
		},
		{
			name:     "should convert tag names to lowercase",
			input:    "Ability.Damage.Fire",
			expected: "ability.damage.fire",
		},
		{
			name:     "should remove trailing dots",
			input:    "ability.damage.fire.",
			expected: "ability.damage.fire",
		},
		{
			name:     "should remove special characters",
			input:    "ability.dam@#$%^&*age.fire.",
			expected: "ability.damage.fire",
		},
		{
			name:     "should remove unicode characters",
			input:    "ability.🔥",
			expected: "ability",
		},
		{
			name:     "should remove extra whitespaces",
			input:    "ability  .   damage .  fire",
			expected: "ability.damage.fire",
		},
		{
			name:     "should remove tabs",
			input:    "ability.\tdamage.fire",
			expected: "ability.damage.fire",
		},
		{
			name:     "should remove trailing or leading whitespace",
			input:    "  ability.damage.fire  ",
			expected: "ability.damage.fire",
		},
		{
			name:     "should normalize consecutive dots",
			input:    "ability.damage..fire",
			expected: "ability.damage.fire",
		},
		{
			name:     "should remove leading dots",
			input:    ".ability.damage.fire",
			expected: "ability.damage.fire",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := tag.FromString(tt.input)
			assert.Equal(t, tt.expected, tag.String())
		})
	}
}

func TestTag_MatchExact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tag1     string
		tag2     string
		expected bool
	}{
		{
			name:     "exact matches",
			tag1:     "Ability.Damage.Fire",
			tag2:     "Ability.Damage.Fire",
			expected: true,
		},
		{
			name:     "non-exact matches",
			tag1:     "Ability.Damage",
			tag2:     "Ability.Damage.Fire",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag1 := tag.FromString(tt.tag1)
			tag2 := tag.FromString(tt.tag2)
			assert.Equal(t, tt.expected, tag1.MatchExact(tag2))
		})
	}
}

func TestTag_Match(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tag1     string
		tag2     string
		expected bool
	}{
		{
			name:     "matches more general tag",
			tag1:     "Ability.Damage.Fire",
			tag2:     "Ability.Damage",
			expected: true,
		},
		{
			name:     "does not match more specific tag",
			tag1:     "Ability.Damage",
			tag2:     "Ability.Damage.Fire",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag1 := tag.FromString(tt.tag1)
			tag2 := tag.FromString(tt.tag2)
			assert.Equal(t, tt.expected, tag1.Match(tag2))
		})
	}
}

func TestTag_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "returns the string value",
			input:    "Ability.Damage.Fire",
			expected: "ability.damage.fire",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := tag.FromString(tt.input)
			assert.Equal(t, tt.expected, tag.String())
		})
	}
}
