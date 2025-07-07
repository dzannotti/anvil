package prettyprint

import (
	"strings"
)

// Tree drawing constants
const (
	TreeFork       = "├─ "
	TreeEnd        = "└─ "
	TreeEndCircle  = "└─○"
	TreeVertical   = "│  "
	TreeVertical1  = "│ "
	TreeBranch     = " ├─ "
	TreeBranchEnd  = " └─ "
	TreeSpace      = "    "
	TreeContinue   = " │   "
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

// AddIndentedBlock adds a text block with the same formatting as indent()
func (tb *TreeBuilder) AddIndentedBlock(text string) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return
	}

	spacing := strings.Repeat(TreeVertical, tb.depth)
	
	// First line gets the fork
	tb.lines = append(tb.lines, spacing+TreeFork+lines[0])
	
	// Subsequent lines get vertical continuation
	for i := 1; i < len(lines); i++ {
		if lines[i] != "" {
			tb.lines = append(tb.lines, spacing+TreeVertical+lines[i])
		}
	}
}

// Helper function to indent existing text block
func indentBlock(text string, depth int) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	spacing := strings.Repeat(TreeVertical, depth)
	lines[0] = spacing + TreeFork + lines[0]
	for i := 1; i < len(lines); i++ {
		lines[i] = spacing + TreeVertical + lines[i]
	}

	return strings.Join(lines, "\n")
}