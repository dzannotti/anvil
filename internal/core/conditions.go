package core

import (
	"anvil/internal/tag"
	"slices"
)

// This map may seems useless but we might end up having multiple poison conditions and we only want to remove one
type Conditions struct {
	Sources map[tag.Tag][]*Effect
}

func (c *Conditions) init() {
	if c.Sources == nil {
		c.Sources = make(map[tag.Tag][]*Effect)
	}
}

func (c *Conditions) Has(t tag.Tag, src *Effect) bool {
	c.init()
	if src == nil {
		return len(c.Sources[t]) > 0
	}
	return slices.Contains(c.Sources[t], src)
}

func (c *Conditions) Match(t tag.Tag) bool {
	for tag := range c.Sources {
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
	c.Sources[t] = append(c.Sources[t], src)
}

func (c *Conditions) Remove(t tag.Tag, src *Effect) bool {
	c.init()
	if src == nil {
		return c.removeAll(t)
	}
	return c.removeSpecific(t, src)
}

func (c *Conditions) removeSpecific(t tag.Tag, src *Effect) bool {
	before := len(c.Sources[t])
	c.Sources[t] = slices.DeleteFunc(c.Sources[t], func(fx *Effect) bool { return src == fx })
	after := len(c.Sources[t])
	return after < before
}

func (c *Conditions) removeAll(t tag.Tag) bool {
	before := len(c.Sources[t])
	// Clear slice while preserving capacity
	c.Sources[t] = c.Sources[t][:0]
	after := len(c.Sources[t])
	return after < before
}
