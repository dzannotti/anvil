package prettyprint

import (
	"anvil/internal/core/tags"
	"anvil/internal/tag"
	"fmt"
	"strings"

	"github.com/adam-lavrik/go-imath/ix"
)

func printValue(value int, first bool) string {
	if first {
		return fmt.Sprintf("= %d", value)
	}
	if value > 0 {
		return fmt.Sprintf("+ %d", value)
	}
	return fmt.Sprintf("- %d", ix.Abs(value))
}

func printTerm(term Term, indent string, last, first bool) []string {
	branch := " ├─ "
	if last {
		branch = " └─ "
	}
	result := make([]string, 0)
	value := printValue(term.Value, first)
	source := strings.Builder{}
	source.WriteString(term.Source)
	var formulaPart string

	if strings.Contains(string(term.Type), "dice") {
		formula := fmt.Sprintf("%dd%d", term.Times, term.Sides)
		if strings.Contains(string(term.Type), "20") {
			formattedAdv := make([]string, 0)
			advDisIndent := indent
			if last {
				advDisIndent += "    "
			} else {
				advDisIndent += "│   "
			}

			for idx, advSource := range term.HasAdvantage {
				isLastAdv := idx == len(term.HasAdvantage)+len(term.HasDisadvantage)-1
				advBranch := "├─ "
				if isLastAdv {
					advBranch = "└─ "
				}
				formattedAdv = append(formattedAdv, fmt.Sprintf("\n%s%sAdvantage: %s", advDisIndent, advBranch, advSource))
			}

			for idx, disSource := range term.HasDisadvantage {
				isLastDis := idx == len(term.HasDisadvantage)-1
				disBranch := "├─ "
				if isLastDis {
					disBranch = "└─ "
				}
				formattedAdv = append(formattedAdv, fmt.Sprintf("\n%s%sDisadvantage: %s", advDisIndent, disBranch, disSource))
			}

			if len(formattedAdv) > 0 {
				source.WriteString(strings.Join(formattedAdv, ""))
			}
		}

		if len(term.Values) > 1 {
			rolls := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(term.Values)), ", "), "[]")
			formulaPart = fmt.Sprintf(" (%s: %s)", formula, rolls)
		} else {
			formulaPart = fmt.Sprintf(" (%s)", formula)
		}
	}

	if len(term.Tags.Tags) > 0 && len(indent) == 0 {
		termTags := strings.Builder{}
		for _, t := range term.Tags.Tags {
			termTags.WriteString(tags.ToReadableTag(tag.FromString(t)))
		}
		source.WriteString(fmt.Sprintf(" (%s)", termTags.String()))
	}

	result = append(result, fmt.Sprintf("%s%s%s%s %s", indent, branch, value, formulaPart, source.String()))

	if len(term.Terms) > 0 {
		newIndent := indent
		if last {
			newIndent += "    "
		} else {
			newIndent += "│   "
		}
		result = append(result, printTerms(term.Terms, newIndent)...)
	}
	return result
}

func printTerms(terms []Term, indent string) []string {
	lines := make([]string, 0)
	for i, term := range terms {
		last := i == len(terms)-1
		first := i == 0
		lines = append(lines, printTerm(term, indent, last, first)...)
	}
	return lines
}

func printExpression(exp Expression) string {
	lines := make([]string, 1)
	lines[0] = fmt.Sprintf("%d", exp.Value)
	lines = append(lines, printTerms(exp.Terms, "")...)
	return strings.Join(lines, "\n")
}
