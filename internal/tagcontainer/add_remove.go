package tagcontainer

import (
	"anvil/internal/tag"
)

func (tc *TagContainer) Add(other TagContainer) {
	for _, t := range other.tags {
		tc.AddTag(t)
	}
}

func (tc *TagContainer) AddTag(newTag tag.Tag) {
	if tc.HasTag(newTag) {
		return
	}
	tc.tags = append(tc.tags, newTag)
}

func (tc *TagContainer) RemoveTag(target tag.Tag) {
	newTags := make([]tag.Tag, 0, len(tc.tags))
	for _, existing := range tc.tags {
		if !existing.MatchExact(target) {
			newTags = append(newTags, existing)
		}
	}
	tc.tags = newTags
}
