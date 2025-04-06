package core

import (
	"anvil/internal/tag"
	"slices"
)

// This map may seems useless but we might end up having multiple poison conditions and we only want to remove one
type Conditions struct {
	Conditions map[tag.Tag][]*Effect
}

func (c *Conditions) init() {
	if c.Conditions == nil {
		c.Conditions = make(map[tag.Tag][]*Effect, 0)
	}
}

func (c *Conditions) Has(t tag.Tag, src *Effect) bool {
	c.init()
	if src == nil {
		return len(c.Conditions[t]) > 0
	}
	return slices.Contains(c.Conditions[t], src)
}

func (c *Conditions) Match(t tag.Tag) bool {
	tags := make([]tag.Tag, 0)
	for tag := range c.Conditions {
		tags = append(tags, tag)
	}
	container := tag.ContainerFromTag(tags...)
	return container.MatchTag(t)
}

func (c *Conditions) Add(t tag.Tag, src *Effect) {
	c.init()
	if src == nil {
		return
	}
	c.Conditions[t] = append(c.Conditions[t], src)
}

func (c *Conditions) Remove(t tag.Tag, src *Effect) bool {
	c.init()
	if src == nil {
		return c.removeAll(t)
	}
	return c.removeSpecific(t, src)
}

func (c *Conditions) removeSpecific(t tag.Tag, src *Effect) bool {
	before := len(c.Conditions[t])
	c.Conditions[t] = slices.DeleteFunc(c.Conditions[t], func(fx *Effect) bool { return src == fx })
	after := len(c.Conditions[t])
	return after < before
}

func (c *Conditions) removeAll(t tag.Tag) bool {
	before := len(c.Conditions[t])
	c.Conditions[t] = make([]*Effect, 0)
	after := len(c.Conditions[t])
	return after < before
}
