package tag

import (
	"slices"
	"strings"
)

type Container struct {
	tags []Tag
}

func NewContainerFromString(values ...string) Container {
	c := Container{}
	for _, v := range values {
		c.AddTag(New(v))
	}

	return c
}

func NewContainer(tag ...Tag) Container {
	c := Container{}
	c.AddTag(tag...)
	return c
}

func NewContainerFromContainer(other Container) Container {
	return NewContainer(other.tags...)
}

func (c *Container) AsStrings() []string {
	if len(c.tags) == 0 {
		return []string{}
	}

	out := make([]string, len(c.tags))
	for i, t := range c.tags {
		out[i] = t.AsString()
	}

	return out
}

func (c *Container) AddTag(tags ...Tag) {
	for _, t := range tags {
		if c.HasTag(t) {
			continue
		}

		c.tags = append(c.tags, t)
	}
}

func (c *Container) Add(container ...Container) {
	for _, other := range container {
		c.AddTag(other.tags...)
	}
}

func (c *Container) RemoveTag(tags ...Tag) {
	removeSet := make(map[Tag]struct{}, len(tags))
	for _, t := range tags {
		removeSet[t] = struct{}{}
	}

	newTags := c.tags[:0]
	for _, t := range c.tags {
		if _, shouldRemove := removeSet[t]; !shouldRemove {
			newTags = append(newTags, t)
		}
	}

	c.tags = newTags
}

func (c *Container) HasTag(other Tag) bool {
	for _, existing := range c.tags {
		if existing.MatchExact(other) {
			return true
		}
	}

	return false
}

func (c *Container) HasAny(other Container) bool {
	for _, t := range other.tags {
		if c.HasTag(t) {
			return true
		}
	}

	return false
}

func (c *Container) HasAll(other Container) bool {
	for _, t := range other.tags {
		if !c.HasTag(t) {
			return false
		}
	}

	return true
}

func (c *Container) MatchTag(tag Tag) bool {
	for _, t := range c.tags {
		if t.Match(tag) {
			return true
		}
	}

	return false
}

func (c *Container) MatchAny(other Container) bool {
	for _, t := range other.tags {
		if c.MatchTag(t) {
			return true
		}
	}

	return false
}

func (c *Container) MatchAll(other Container) bool {
	for _, t := range other.tags {
		if !c.MatchTag(t) {
			return false
		}
	}

	return true
}

func (c *Container) IsEmpty() bool {
	return len(c.tags) == 0
}

func (c *Container) ID() string {
	strs := c.AsStrings()
	slices.Sort(strs)
	return strings.Join(strs, "-")
}

func (c *Container) Len() int {
	return len(c.tags)
}

func (c *Container) Clone() Container {
	return NewContainer(c.tags...)
}
