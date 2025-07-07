package prettyprint

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeBuilder_BasicStructure(t *testing.T) {
	tb := NewTreeBuilder()
	tb.AddRawLine("Root")
	tb.AddLine("Child 1")
	tb.AddLine("Child 2")

	result := tb.String()
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 3)
	assert.Equal(t, "Root", lines[0])
	assert.Equal(t, "├─ Child 1", lines[1])
	assert.Equal(t, "├─ Child 2", lines[2])
}

func TestTreeBuilder_WithIndent(t *testing.T) {
	tb := NewTreeBuilder()
	tb.AddRawLine("Root")
	tb.WithIndent(func() {
		tb.AddLine("Nested Child")
	})
	tb.AddLine("Sibling")

	result := tb.String()
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 3)
	assert.Equal(t, "Root", lines[0])
	assert.Equal(t, "│  ├─ Nested Child", lines[1])
	assert.Equal(t, "├─ Sibling", lines[2])
}

func TestTreeBuilder_AddBranch(t *testing.T) {
	tb := NewTreeBuilder()
	tb.AddBranch("First", false)
	tb.AddBranch("Last", true)

	result := tb.String()
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 2)
	assert.Equal(t, "├─ First", lines[0])
	assert.Equal(t, "└─ Last", lines[1])
}

func TestTreeBuilder_AddIndentedBlock(t *testing.T) {
	tb := NewTreeBuilder()
	tb.AddIndentedBlock("Line 1\nLine 2\nLine 3")

	result := tb.String()
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 3)
	assert.Equal(t, "├─ Line 1", lines[0])
	assert.Equal(t, "│  Line 2", lines[1])
	assert.Equal(t, "│  Line 3", lines[2])
}

func TestTreeBuilder_AddIndentedBlock_EmptyLines(t *testing.T) {
	tb := NewTreeBuilder()
	tb.AddIndentedBlock("Line 1\n\nLine 3")

	result := tb.String()
	lines := strings.Split(result, "\n")

	// Empty lines should be skipped
	assert.Len(t, lines, 2)
	assert.Equal(t, "├─ Line 1", lines[0])
	assert.Equal(t, "│  Line 3", lines[1])
}

func TestIndentBlock(t *testing.T) {
	text := "First line\nSecond line\nThird line"
	result := indentBlock(text, 1)

	expected := "│  ├─ First line\n│  │  Second line\n│  │  Third line"
	assert.Equal(t, expected, result)
}

func TestIndentBlock_EmptyText(t *testing.T) {
	result := indentBlock("", 1)
	assert.Equal(t, "│  ├─ ", result)
}

func TestIndentBlock_ZeroDepth(t *testing.T) {
	text := "Single line"
	result := indentBlock(text, 0)
	assert.Equal(t, "├─ Single line", result)
}

func TestGetChildIndent(t *testing.T) {
	tests := []struct {
		name     string
		indent   string
		last     bool
		expected string
	}{
		{
			name:     "not last",
			indent:   "│  ",
			last:     false,
			expected: "│   │   ",
		},
		{
			name:     "last",
			indent:   "│  ",
			last:     true,
			expected: "│      ",
		},
		{
			name:     "root not last",
			indent:   "",
			last:     false,
			expected: " │   ",
		},
		{
			name:     "root last",
			indent:   "",
			last:     true,
			expected: "    ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getChildIndent(tt.indent, tt.last)
			assert.Equal(t, tt.expected, result)
		})
	}
}