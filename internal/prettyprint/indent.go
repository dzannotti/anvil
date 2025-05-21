package prettyprint

import (
	"strings"
)

func indent(text string) string {
	depth := 0
	spacing := strings.Repeat("│  ", max(0, depth))
	lines := strings.Split(text, "\n")
	lines[0] = spacing + "├─" + lines[0]
	for i, line := range lines[1:] {
		lines[i+1] = spacing + "│ " + line
	}
	return strings.Join(lines, "\n")
}
