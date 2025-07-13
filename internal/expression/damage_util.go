package expression

import (
	"fmt"
	"anvil/internal/tag"
)

func groupComponentsByTags(components []Component) [][]Component {
	groups := make([][]Component, 0, len(components))
	used := make([]bool, len(components))

	for i, comp := range components {
		if used[i] {
			continue
		}
		group := []Component{comp}
		used[i] = true
		compTags := resolveComponentTags(comp, components)

		for j := i + 1; j < len(components); j++ {
			if used[j] {
				continue
			}
			otherTags := resolveComponentTags(components[j], components)
			if compTags.ID() == otherTags.ID() {
				group = append(group, components[j])
				used[j] = true
			}
		}
		groups = append(groups, group)
	}
	return groups
}

func resolveComponentTags(component Component, components []Component) tag.Container {
	componentTags := component.Tags()
	if componentTags.HasTag(Primary) && len(components) > 0 {
		return components[0].Tags()
	}

	return componentTags
}

func resolveGroupTags(group []Component, components []Component) tag.Container {
	if len(group) == 0 {
		return tag.Container{}
	}

	return resolveComponentTags(group[0], components)
}

func buildGroupSource(group []Component) string {
	if len(group) == 0 {
		return "empty group"
	}

	if len(group) == 1 {
		return group[0].Source()
	}

	if primarySource := findPrimarySource(group); primarySource != "" {
		return primarySource
	}

	return buildMultiComponentSource(group)
}

func findPrimarySource(group []Component) string {
	primaryComponents := 0
	var primarySource string

	for _, comp := range group {
		tags := comp.Tags()
		if tags.HasTag(Primary) || tags.IsEmpty() {
			primaryComponents++
		} else if primarySource == "" {
			primarySource = comp.Source()
		}
	}

	if primarySource != "" && primaryComponents > 0 {
		return primarySource
	}
	return ""
}

func buildMultiComponentSource(group []Component) string {
	firstSource := group[0].Source()
	for _, comp := range group[1:] {
		if comp.Source() != firstSource {
			return fmt.Sprintf("grouped damage (%d sources)", len(group))
		}
	}
	return firstSource
}
