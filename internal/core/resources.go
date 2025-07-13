package core

import (
	"maps"

	"anvil/internal/core/tags"
	"anvil/internal/loader"
	"anvil/internal/mathi"
	"anvil/internal/tag"
)

type Resources struct {
	Current map[tag.Tag]int
	Max     map[tag.Tag]int
}

func NewResourcesFromDefinition(def loader.ResourcesDefinition) Resources {
	resources := Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed: def.WalkSpeed,
	}}
	
	optionalResources := map[tag.Tag]int{
		tags.FlySpeed:    def.FlySpeed,
		tags.SwimSpeed:   def.SwimSpeed,
		tags.SpellSlot1:  def.SpellSlot1,
		tags.SpellSlot2:  def.SpellSlot2,
		tags.SpellSlot3:  def.SpellSlot3,
		tags.SpellSlot4:  def.SpellSlot4,
		tags.SpellSlot5:  def.SpellSlot5,
		tags.SpellSlot6:  def.SpellSlot6,
		tags.SpellSlot7:  def.SpellSlot7,
		tags.SpellSlot8:  def.SpellSlot8,
		tags.SpellSlot9:  def.SpellSlot9,
	}
	
	for resource, value := range optionalResources {
		if value > 0 {
			resources.Max[resource] = value
		}
	}
	
	return resources
}

func (r *Resources) init() {
	if r.Current == nil {
		r.Current = make(map[tag.Tag]int)
	}
}

func (r *Resources) Reset() {
	r.init()
	r.Current[tags.Action] = 1
	r.Current[tags.BonusAction] = 1
	r.Current[tags.Reaction] = 1
	r.Current[tags.UsedSpeed] = 0
}

func (r *Resources) LongRest() {
	r.init()
	maps.Copy(r.Current, r.Max)
	r.Reset()
}

func (r Resources) CanUse(t tag.Tag, v int) bool {
	r.init()
	return r.Current[t] >= v
}

func (r Resources) CanAfford(c map[tag.Tag]int) bool {
	r.init()
	for t, v := range c {
		if r.Current[t] < v {
			return false
		}
	}
	return true
}

func (r Resources) Consume(t tag.Tag, v int) {
	if t.Match(tags.Speed) {
		r.Current[tags.UsedSpeed] += v
		return
	}
	r.Current[t] -= v
}

func (r Resources) Remaining(t tag.Tag) int {
	r.init()
	if t.Match(tags.Speed) {
		return r.remainingSpeed(t)
	}
	return r.Current[t]
}

func (r Resources) remainingSpeed(t tag.Tag) int {
	r.init()
	m := r.Max[t]
	total := r.maxSpeed() - r.Current[tags.UsedSpeed]
	remaining := mathi.Min(m-r.Current[tags.UsedSpeed], total)
	if remaining <= 0 {
		return 0
	}
	return remaining
}

func (r Resources) maxSpeed() int {
	r.init()
	m := 0
	for k, v := range r.Max {
		if !k.Match(tags.Speed) {
			continue
		}
		if v > m {
			m = v
		}
	}
	return m
}
