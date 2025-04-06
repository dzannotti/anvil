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
	for tag := range c.Conditions {
		if tag.Match(t) {
			return true
		}
	}
	return false
}

func (c *Conditions) Add(t tag.Tag, src *Effect) {
	c.init()
	if src == nil {
		panic("Attempted to add a condition with no source")
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
	// Clear slice while preserving capacity
	c.Conditions[t] = c.Conditions[t][:0]
	after := len(c.Conditions[t])
	return after < before
}
