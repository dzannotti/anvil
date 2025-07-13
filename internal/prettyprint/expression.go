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

func formatDiceRolls(component expression.Component) string {
	var values []int
	var times, sides int

	// Type assert to get dice-specific methods
	switch c := component.(type) {
	case interface {
		Values() []int
		Times() int
		Sides() int
	}:
		values = c.Values()
		times = c.Times()
		sides = c.Sides()
	default:
		return ""
	}

	if len(values) <= 1 {
		return ""
	}

	rolls := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(values)), ", "), "[]")
	formula := fmt.Sprintf("%dd%d", times, sides)
	return fmt.Sprintf(" (%s: %s)", formula, rolls)
}

func formatDice(component expression.Component) string {
	if component.Kind() != expression.ComponentKindDice && component.Kind() != expression.ComponentKindD20 {
		return ""
	}

	// Type assert to get dice-specific methods
	switch c := component.(type) {
	case interface {
		Values() []int
		Times() int
		Sides() int
	}:
		if len(c.Values()) <= 1 {
			return fmt.Sprintf(" (%dd%d)", c.Times(), c.Sides())
		}
		return formatDiceRolls(component)
	case interface{ Values() []int }:
		// D20 component - just show d20
		if len(c.Values()) <= 1 {
			return " (1d20)"
		}
		return formatDiceRolls(component)
	default:
		return ""
	}
}

func formatBranch(indent string, last bool) string {
	if last {
		return indent + TreeSpace
	}
	return indent + TreeContinue
}

func formatAdvantageDisadvantage(component expression.Component, indent string, last bool) []string {
	var advantages, disadvantages []string

	// Type assert to get D20-specific methods
	if d20, ok := component.(interface {
		Advantage() []string
		Disadvantage() []string
	}); ok {
		advantages = d20.Advantage()
		disadvantages = d20.Disadvantage()
	} else {
		return nil
	}

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

func formatTags(component expression.Component) string {
	compTags := component.Tags()
	if compTags.IsEmpty() {
		return ""
	}
	componentTags := make([]string, len(compTags.AsStrings()))
	for i, t := range compTags.AsStrings() {
		componentTags[i] = tags.ToReadable(tag.FromString(t))
	}
	return fmt.Sprintf(" (%s)", strings.Join(componentTags, ", "))
}

func buildComponentSource(component expression.Component, indent string, last bool) string {
	source := strings.Builder{}
	source.WriteString(component.Source())

	if component.Kind() == expression.ComponentKindD20 {
		advDisadv := formatAdvantageDisadvantage(component, indent, last)
		if len(advDisadv) > 0 {
			source.WriteString(strings.Join(advDisadv, ""))
		}
	}

	compTags := component.Tags()
	if !compTags.IsEmpty() && len(indent) == 0 {
		source.WriteString(formatTags(component))
	}

	return source.String()
}

func getChildIndent(indent string, last bool) string {
	if last {
		return indent + TreeSpace
	}
	return indent + TreeContinue
}

func printComponent(component expression.Component, indent string, last, first bool) []string {
	branch := TreeBranch
	if last {
		branch = TreeBranchEnd
	}

	result := make([]string, 0)
	value := printValue(component.Value(), first)
	source := buildComponentSource(component, indent, last)

	// Build the main line
	result = append(result, fmt.Sprintf("%s%s%s%s %s", indent, branch, value, formatDice(component), source))

	// Add child components if any
	if len(component.Components()) > 0 {
		childIndent := getChildIndent(indent, last)
		result = append(result, printComponents(component.Components(), childIndent)...)
	}

	return result
}

func printComponents(components []expression.Component, indent string) []string {
	lines := make([]string, 0)
	for i, component := range components {
		last := i == len(components)-1
		first := i == 0
		lines = append(lines, printComponent(component, indent, last, first)...)
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

	// Add components using the existing logic
	componentLines := printComponents(exp.Components, "")
	for _, line := range componentLines {
		tb.AddRawLine(line)
	}

	return tb.String()
}
