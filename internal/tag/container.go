package tag

import (
	"strings"
)

type Container struct {
	tags []Tag
}

func NewContainer() Container {
	return Container{tags: []Tag{}}
}

func (tc Container) Strings() []string {
	tags := make([]string, len(tc.tags))
	for i, tag := range tc.tags {
		tags[i] = tag.String()
	}
	return tags
}

func (tc Container) ID() string {
	return strings.Join(tc.Strings(), "-")
}

func (tc Container) Clone() Container {
	return ContainerFromTags(tc.tags)
}

func (tc Container) IsEmpty() bool {
	return len(tc.tags) == 0
}
