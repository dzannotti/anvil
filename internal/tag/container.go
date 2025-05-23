package tag

import (
	"slices"
	"strings"
)

type Container struct {
	tags []Tag
}

func ContainerFromString(values ...string) Container {
	c := Container{}
	for _, v := range values {
		c.AddTag(FromString(v))
	}
	return c
}

func ContainerFromTag(tag ...Tag) Container {
	c := Container{}
	c.AddTag(tag...)
	return c
}

func ContainerFromContainer(other Container) Container {
	return ContainerFromTag(other.tags...)
}

func (c *Container) AsStrings() []string {
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

func (c *Container) HasAnyTag(other Container) bool {
	for _, t := range other.tags {
		if c.HasTag(t) {
			return true
		}
	}
	return false
}

func (c *Container) HasAllTag(other Container) bool {
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

func (c *Container) MatchAnyTag(other Container) bool {
	for _, t := range other.tags {
		if c.MatchTag(t) {
			return true
		}
	}
	return false
}

func (c *Container) MatchAllTag(other Container) bool {
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
