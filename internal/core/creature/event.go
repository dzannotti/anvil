package creature

import "anvil/internal/core/event"

func ToEvent(c *Creature) event.Creature {
	return event.Creature{
		Name:      c.Name(),
		FactionID: c.FactionID(),
		HitPoints: c.HitPoints(),
	}
}

func NewTakeDamageEvent(c *Creature, d int) *event.TakeDamage {
	return &event.TakeDamage{
		Creature: ToEvent(c),
		Damage:   d,
	}
}

func NewDeathEvent(c *Creature) *event.Death {
	return &event.Death{
		Creature: ToEvent(c),
	}
}

func NewUseActionEvent(a string, source *Creature, target *Creature) *event.UseAction {
	return &event.UseAction{
		Action: event.Action{
			Name: a,
		},
		Source: ToEvent(source),
	}
}
