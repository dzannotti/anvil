package tagcontainer

import "anvil/internal/tag"

func (tc TagContainer) MatchTag(target tag.Tag) bool {
	for _, existing := range tc.tags {
		if existing.Match(target) {
			return true
		}
	}
	return false
}

func (tc TagContainer) MatchAnyTag(other TagContainer) bool {
	for _, tag := range other.tags {
		if tc.MatchTag(tag) {
			return true
		}
	}
	return false
}

func (tc TagContainer) MatchAllTag(other TagContainer) bool {
	for _, tag := range other.tags {
		if !tc.MatchTag(tag) {
			return false
		}
	}
	return true
}
