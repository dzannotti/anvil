package prettyprint

import (
	"fmt"
	"strings"

	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/mathi"
	"anvil/internal/tag"
)

func printValue(value int, first bool) string {
	if first {
		return fmt.Sprintf("= %d", value)
	}
	if value > 0 {
		return fmt.Sprintf("+ %d", value)
	}
	return fmt.Sprintf("- %d", mathi.Abs(value))
}

func formatDiceRolls(term expression.Term) string {
	if len(term.Values) <= 1 {
		return ""
	}
	rolls := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(term.Values)), ", "), "[]")
	formula := fmt.Sprintf("%dd%d", term.Times, term.Sides)
	return fmt.Sprintf(" (%s: %s)", formula, rolls)
}

func formatDice(term expression.Term) string {
	if !strings.Contains(string(term.Type), "dice") {
		return ""
	}
	if len(term.Values) <= 1 {
		return fmt.Sprintf(" (%dd%d)", term.Times, term.Sides)
	}
	return formatDiceRolls(term)
}

func formatBranch(indent string, last bool) string {
	if last {
		return indent + "    "
	}
	return indent + "│   "
}

func formatAdvantage(term expression.Term, indent string, last bool) []string {
	if len(term.HasAdvantage) == 0 {
		return nil
	}

	formatted := make([]string, 0)
	indent = formatBranch(indent, last)
	totalItems := len(term.HasAdvantage) + len(term.HasDisadvantage)

	for idx, source := range term.HasAdvantage {
		isLast := idx == totalItems-1
		branch := "├─ "
		if isLast {
			branch = "└─ "
		}
		formatted = append(formatted, fmt.Sprintf("\n%s%sAdvantage: %s", indent, branch, source))
	}

	return formatted
}

func formatDisadvantage(term expression.Term, indent string, last bool) []string {
	if len(term.HasDisadvantage) == 0 {
		return nil
	}

	formatted := make([]string, 0)
	indent = formatBranch(indent, last)

	for idx, source := range term.HasDisadvantage {
		isLast := idx == len(term.HasDisadvantage)-1
		branch := "├─ "
		if isLast {
			branch = "└─ "
		}
		formatted = append(formatted, fmt.Sprintf("\n%s%sDisadvantage: %s", indent, branch, source))
	}

	return formatted
}

func formatTags(term expression.Term) string {
	if term.Tags.IsEmpty() {
		return ""
	}
	termTags := make([]string, len(term.Tags.Strings()))
	for i, t := range term.Tags.Strings() {
		termTags[i] = tags.ToReadable(tag.FromString(t))
	}
	return fmt.Sprintf(" (%s)", strings.Join(termTags, ", "))
}

func printTerm(term expression.Term, indent string, last, first bool) []string {
	branch := " ├─ "
	if last {
		branch = " └─ "
	}
	result := make([]string, 0)
	value := printValue(term.Value, first)
	source := strings.Builder{}
	source.WriteString(term.Source)

	if strings.Contains(string(term.Type), "dice") {
		advantages := formatAdvantage(term, indent, last)
		disadvantages := formatDisadvantage(term, indent, last)

		if len(advantages) > 0 {
			source.WriteString(strings.Join(advantages, ""))
		}
		if len(disadvantages) > 0 {
			source.WriteString(strings.Join(disadvantages, ""))
		}
	}

	if !term.Tags.IsEmpty() && len(indent) == 0 {
		source.WriteString(formatTags(term))
	}

	result = append(result, fmt.Sprintf("%s%s%s%s %s", indent, branch, value, formatDice(term), source.String()))

	if len(term.Terms) > 0 {
		newIndent := indent
		if last {
			newIndent += "    "
		} else {
			newIndent += " │   "
		}
		result = append(result, printTerms(term.Terms, newIndent)...)
	}
	return result
}

func printTerms(terms []expression.Term, indent string) []string {
	lines := make([]string, 0)
	for i, term := range terms {
		last := i == len(terms)-1
		first := i == 0
		lines = append(lines, printTerm(term, indent, last, first)...)
	}
	return lines
}

func printExpression(exp *expression.Expression, start ...bool) string {
	lines := make([]string, 1)
	space := ""
	if len(start) > 0 && start[0] {
		space = " "
	}
	lines[0] = fmt.Sprintf("%s%d", space, exp.Value)
	lines = append(lines, printTerms(exp.Terms, "")...)
	return strings.Join(lines, "\n")
}
