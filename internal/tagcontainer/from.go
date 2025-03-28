package tagcontainer

import "anvil/internal/tag"

func New() *TagContainer {
	return &TagContainer{tags: []tag.Tag{}}
}

func FromString(value string) *TagContainer {
	return &TagContainer{
		tags: []tag.Tag{tag.FromString(value)},
	}
}

func FromStrings(values []string) *TagContainer {
	tags := make([]tag.Tag, len(values))
	for i, value := range values {
		tags[i] = tag.FromString(value)
	}
	return &TagContainer{tags: tags}
}

func FromTag(value tag.Tag) *TagContainer {
	return &TagContainer{
		tags: []tag.Tag{value},
	}
}

func FromTags(values []tag.Tag) *TagContainer {
	tags := make([]tag.Tag, len(values))
	copy(tags, values)
	return &TagContainer{tags: tags}
}

func FromContainer(value TagContainer) *TagContainer {
	return FromStrings(value.Strings())
}
