package prettyprint

import (
	"strings"
)

// Tree drawing constants
const (
	TreeFork      = "├─ "
	TreeEnd       = "└─ "
	TreeEndCircle = "└─○"
	TreeVertical  = "│  "
	TreeBranch    = " ├─ "
	TreeBranchEnd = " └─ "
	TreeSpace     = "    "
	TreeContinue  = " │   "
)

type TreeBuilder struct {
	lines []string
	depth int
}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		lines: make([]string, 0),
		depth: 0,
	}
}

func (tb *TreeBuilder) AddLine(content string) {
	spacing := strings.Repeat(TreeVertical, tb.depth)
	tb.lines = append(tb.lines, spacing+TreeFork+content)
}

func (tb *TreeBuilder) AddRawLine(content string) {
	spacing := strings.Repeat(TreeVertical, tb.depth)
	tb.lines = append(tb.lines, spacing+content)
}

func (tb *TreeBuilder) AddBranch(content string, last bool) {
	spacing := strings.Repeat(TreeVertical, tb.depth)
	branch := TreeFork
	if last {
		branch = TreeEnd
	}
	tb.lines = append(tb.lines, spacing+branch+content)
}

func (tb *TreeBuilder) AddEnd() {
	if tb.depth > 0 {
		spacing := strings.Repeat(TreeVertical, tb.depth-1)
		tb.lines = append(tb.lines, spacing+TreeEndCircle)
	}
}

func (tb *TreeBuilder) Indent() {
	tb.depth++
}

func (tb *TreeBuilder) Outdent() {
	if tb.depth > 0 {
		tb.depth--
	}
}

func (tb *TreeBuilder) String() string {
	return strings.Join(tb.lines, "\n")
}

func (tb *TreeBuilder) WithIndent(fn func()) {
	tb.Indent()
	fn()
	tb.Outdent()
}

// processIndentedLines applies tree indentation to text lines, skipping empty lines
func processIndentedLines(lines []string, depth int) []string {
	if len(lines) == 0 {
		return lines
	}

	spacing := strings.Repeat(TreeVertical, depth)
	result := make([]string, 0, len(lines))

	// First non-empty line gets the fork
	result = append(result, spacing+TreeFork+lines[0])

	// Subsequent lines get vertical continuation, skip empty lines
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		result = append(result, spacing+TreeVertical+lines[i])
	}

	return result
}

// AddIndentedBlock adds a text block with tree indentation
func (tb *TreeBuilder) AddIndentedBlock(text string) {
	lines := strings.Split(text, "\n")
	indentedLines := processIndentedLines(lines, tb.depth)
	tb.lines = append(tb.lines, indentedLines...)
}

// Helper function to indent existing text block
func indentBlock(text string, depth int) string {
	lines := strings.Split(text, "\n")
	indentedLines := processIndentedLines(lines, depth)
	return strings.Join(indentedLines, "\n")
}
