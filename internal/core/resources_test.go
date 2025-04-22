package core

import (
	"testing"

	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func TestResources(t *testing.T) {
	t.Run("Action Economy", func(t *testing.T) {
		t.Run("should start with one action, bonus action, and reaction per turn", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.Action:      1,
					tags.BonusAction: 1,
					tags.Reaction:    1,
				},
			}
			resources.Reset()
			resources.LongRest()
			if resources.Remaining(tags.Action) != 1 {
				t.Errorf("expected 1 action, got %d", resources.Remaining(tags.Action))
			}
			if resources.Remaining(tags.BonusAction) != 1 {
				t.Errorf("expected 1 bonus action, got %d", resources.Remaining(tags.BonusAction))
			}
			if resources.Remaining(tags.Reaction) != 1 {
				t.Errorf("expected 1 reaction, got %d", resources.Remaining(tags.Reaction))
			}
		})

		t.Run("should consume actions when used", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.Action:      1,
					tags.BonusAction: 1,
					tags.Reaction:    1,
				},
			}
			resources.Reset()
			resources.LongRest()
			resources.Consume(tags.Action, 1)
			if resources.Remaining(tags.Action) != 0 {
				t.Errorf("expected 0 actions, got %d", resources.Remaining(tags.Action))
			}
			if resources.CanUse(tags.Action, 1) {
				t.Error("should not be able to use action")
			}
		})

		t.Run("should reset all actions at the start of new turn", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.Action:      1,
					tags.BonusAction: 1,
					tags.Reaction:    1,
				},
			}
			resources.Reset()
			resources.LongRest()
			resources.Consume(tags.Action, 1)
			resources.Consume(tags.BonusAction, 1)
			resources.Consume(tags.Reaction, 1)

			resources.Reset()

			if resources.Remaining(tags.Action) != 1 {
				t.Errorf("expected 1 action after reset, got %d", resources.Remaining(tags.Action))
			}
			if resources.Remaining(tags.BonusAction) != 1 {
				t.Errorf("expected 1 bonus action after reset, got %d", resources.Remaining(tags.BonusAction))
			}
			if resources.Remaining(tags.Reaction) != 1 {
				t.Errorf("expected 1 reaction after reset, got %d", resources.Remaining(tags.Reaction))
			}
		})
	})

	t.Run("Movement Speed", func(t *testing.T) {
		t.Run("should calculate max speed and track remaining movement correctly", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.WalkSpeed: 30,
					tags.SwimSpeed: 20,
					tags.FlySpeed:  40,
				},
			}
			resources.Reset()
			resources.LongRest()
			if maxSpeed := resources.maxSpeed(); maxSpeed != 40 {
				t.Errorf("expected max speed 40, got %d", maxSpeed)
			}

			resources.Consume(tags.WalkSpeed, 10)
			if remaining := resources.Remaining(tags.WalkSpeed); remaining != 20 {
				t.Errorf("expected 20 walk speed remaining, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.SwimSpeed); remaining != 10 {
				t.Errorf("expected 10 swim speed remaining, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.FlySpeed); remaining != 30 {
				t.Errorf("expected 30 fly speed remaining, got %d", remaining)
			}
		})

		t.Run("should not allow movement beyond available speed", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.WalkSpeed: 30,
					tags.SwimSpeed: 20,
					tags.FlySpeed:  40,
				},
			}
			resources.Reset()

			resources.Consume(tags.WalkSpeed, 25)
			if remaining := resources.Remaining(tags.WalkSpeed); remaining != 5 {
				t.Errorf("expected 5 walk speed remaining, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.SwimSpeed); remaining != 0 {
				t.Errorf("expected 0 swim speed remaining, got %d", remaining)
			}
		})

		t.Run("should reset movement speed on new turn", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.WalkSpeed: 30,
					tags.SwimSpeed: 20,
					tags.FlySpeed:  40,
				},
			}
			resources.Reset()

			resources.Consume(tags.WalkSpeed, 15)
			resources.Reset()

			if remaining := resources.Remaining(tags.WalkSpeed); remaining != 30 {
				t.Errorf("expected 30 walk speed after reset, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.SwimSpeed); remaining != 20 {
				t.Errorf("expected 20 swim speed after reset, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.FlySpeed); remaining != 40 {
				t.Errorf("expected 40 fly speed after reset, got %d", remaining)
			}
		})

		t.Run("should handle different movement types independently", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.WalkSpeed: 30,
					tags.SwimSpeed: 20,
					tags.FlySpeed:  40,
				},
			}
			resources.Reset()

			resources.Consume(tags.FlySpeed, 20)
			if remaining := resources.Remaining(tags.FlySpeed); remaining != 20 {
				t.Errorf("expected 20 fly speed remaining, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.WalkSpeed); remaining != 10 {
				t.Errorf("expected 10 walk speed remaining, got %d", remaining)
			}
		})
	})

	t.Run("Custom Resources", func(t *testing.T) {
		t.Run("should handle custom resources", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.LegendaryAction: 3,
					tags.SorceryPoints:   5,
				},
			}

			resources.Reset()
			resources.LongRest()

			if remaining := resources.Remaining(tags.LegendaryAction); remaining != 3 {
				t.Errorf("expected 3 legendary actions, got %d", remaining)
			}
			if remaining := resources.Remaining(tags.SorceryPoints); remaining != 5 {
				t.Errorf("expected 5 sorcery points, got %d", remaining)
			}

			resources.Consume(tags.LegendaryAction, 2)
			if remaining := resources.Remaining(tags.LegendaryAction); remaining != 1 {
				t.Errorf("expected 1 legendary action after consuming 2, got %d", remaining)
			}
		})
	})
}
