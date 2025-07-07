package prettyprint

import (
	"strings"
)

func indent(text string, depth int) string {
	spacing := strings.Repeat("│  ", max(0, depth))
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	lines[0] = spacing + fork + lines[0]
	for i, line := range lines[1:] {
		lines[i+1] = spacing + "│  " + line
	}

	return strings.Join(lines, "\n")
}
