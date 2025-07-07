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
		return indent + TreeSpace
	}
	return indent + TreeContinue
}

func formatAdvantageDisadvantage(term expression.Term, indent string, last bool) []string {
	advantages := term.HasAdvantage
	disadvantages := term.HasDisadvantage

	if len(advantages) == 0 && len(disadvantages) == 0 {
		return nil
	}

	formatted := make([]string, 0)
	baseIndent := formatBranch(indent, last)
	totalItems := len(advantages) + len(disadvantages)

	for idx, source := range advantages {
		isLast := idx == totalItems-1
		branch := TreeFork
		if isLast {
			branch = TreeEnd
		}
		formatted = append(formatted, fmt.Sprintf("\n%s%sAdvantage: %s", baseIndent, branch, source))
	}

	// Add disadvantages
	for idx, source := range disadvantages {
		globalIdx := len(advantages) + idx
		isLast := globalIdx == totalItems-1
		branch := TreeFork
		if isLast {
			branch = TreeEnd
		}
		formatted = append(formatted, fmt.Sprintf("\n%s%sDisadvantage: %s", baseIndent, branch, source))
	}

	return formatted
}

func formatTags(term expression.Term) string {
	if term.Tags.IsEmpty() {
		return ""
	}
	termTags := make([]string, len(term.Tags.AsStrings()))
	for i, t := range term.Tags.AsStrings() {
		termTags[i] = tags.ToReadable(tag.FromString(t))
	}
	return fmt.Sprintf(" (%s)", strings.Join(termTags, ", "))
}

func buildTermSource(term expression.Term, indent string, last bool) string {
	source := strings.Builder{}
	source.WriteString(term.Source)

	if strings.Contains(string(term.Type), "dice") {
		advDisadv := formatAdvantageDisadvantage(term, indent, last)
		if len(advDisadv) > 0 {
			source.WriteString(strings.Join(advDisadv, ""))
		}
	}

	if !term.Tags.IsEmpty() && len(indent) == 0 {
		source.WriteString(formatTags(term))
	}

	return source.String()
}

func getChildIndent(indent string, last bool) string {
	if last {
		return indent + TreeSpace
	}
	return indent + TreeContinue
}

func printTerm(term expression.Term, indent string, last, first bool) []string {
	branch := TreeBranch
	if last {
		branch = TreeBranchEnd
	}

	result := make([]string, 0)
	value := printValue(term.Value, first)
	source := buildTermSource(term, indent, last)

	// Build the main line
	result = append(result, fmt.Sprintf("%s%s%s%s %s", indent, branch, value, formatDice(term), source))

	// Add child terms if any
	if len(term.Terms) > 0 {
		childIndent := getChildIndent(indent, last)
		result = append(result, printTerms(term.Terms, childIndent)...)
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
	tb := NewTreeBuilder()
	space := ""
	if len(start) > 0 && start[0] {
		space = " "
	}

	tb.AddRawLine(fmt.Sprintf("%s%d", space, exp.Value))

	// Add terms using the existing logic
	termLines := printTerms(exp.Terms, "")
	for _, line := range termLines {
		tb.AddRawLine(line)
	}

	return tb.String()
}
