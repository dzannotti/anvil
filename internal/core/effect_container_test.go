package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestContainer struct{}
type UnhandledEvent struct{}

func TestContainer_Add(t *testing.T) {
	t.Run("adds effects in priority order", func(t *testing.T) {
		c := &EffectContainer{}
		e1 := &Effect{Name: "e1", Priority: PriorityNormal}
		e2 := &Effect{Name: "e2", Priority: PriorityEarly}
		e3 := &Effect{Name: "e3", Priority: PriorityLate}

		c.Add(e1)
		c.Add(e2)
		c.Add(e3)

		assert.Equal(t, []*Effect{e2, e1, e3}, c.effects)
	})
}

func TestContainer_Remove(t *testing.T) {
	t.Run("removes effect by name", func(t *testing.T) {
		c := &EffectContainer{}
		e1 := &Effect{Name: "e1"}
		e2 := &Effect{Name: "e2"}

		c.Add(e1)
		c.Add(e2)
		c.Remove(e1)

		assert.Equal(t, []*Effect{e2}, c.effects)
	})

	t.Run("does nothing if effect not found", func(t *testing.T) {
		c := &EffectContainer{}
		e1 := &Effect{Name: "e1"}
		e2 := &Effect{Name: "e2"}
		e3 := &Effect{Name: "e3"}
		c.Add(e1)
		c.Add(e2)
		c.Remove(e3)

		assert.Equal(t, []*Effect{e1, e2}, c.effects)
	})
}

func TestContainer_Evaluate(t *testing.T) {
	t.Run("evaluates effects in priority order", func(t *testing.T) {
		c := &EffectContainer{}
		order := []string{}
		e1 := &Effect{Name: "e1", Priority: PriorityNormal}
		e1.withHandler("TestContainer", func(_ *Effect, _ any) { order = append(order, "e1") })
		e2 := &Effect{Name: "e2", Priority: PriorityEarly}
		e2.withHandler("TestContainer", func(_ *Effect, _ any) { order = append(order, "e2") })
		e3 := &Effect{Name: "e3", Priority: PriorityLate}
		e3.withHandler("TestContainer", func(_ *Effect, _ any) { order = append(order, "e3") })

		c.Add(e1)
		c.Add(e2)
		c.Add(e3)
		c.Evaluate(&TestContainer{})

		assert.Equal(t, []string{"e2", "e1", "e3"}, order)
	})

	t.Run("does nothing for unhandled events", func(t *testing.T) {
		c := &EffectContainer{}
		e1 := &Effect{Name: "e1"}
		e1.withHandler("TestContainer", func(_ *Effect, _ any) { assert.Fail(t, "should not be called") })

		c.Add(e1)
		c.Evaluate(&UnhandledEvent{})
	})
}
