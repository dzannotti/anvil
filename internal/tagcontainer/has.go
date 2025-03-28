package tagcontainer

import "anvil/internal/tag"

func (container TagContainer) HasTag(target tag.Tag) bool {
	for _, t := range container.tags {
		if t.MatchExact(target) {
			return true
		}
	}
	return false
}

func (container TagContainer) HasAnyTag(other TagContainer) bool {
	for _, t := range other.tags {
		if container.HasTag(t) {
			return true
		}
	}
	return false
}

func (container TagContainer) HasAllTag(other TagContainer) bool {
	for _, t := range other.tags {
		if !container.HasTag(t) {
			return false
		}
	}
	return true
}
