package tagcontainer

import (
	"strings"

	"anvil/internal/tag"
)

type TagContainer struct {
	tags []tag.Tag
}

func (tc TagContainer) Strings() []string {
	tags := make([]string, len(tc.tags))
	for i, tag := range tc.tags {
		tags[i] = tag.String()
	}
	return tags
}

func (tc TagContainer) ID() string {
	return strings.Join(tc.Strings(), "-")
}

func (tc TagContainer) Clone() *TagContainer {
	return FromTags(tc.tags)
}
