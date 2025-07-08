package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
			assert.Equal(t, 1, resources.Remaining(tags.Action), "expected 1 action")
			assert.Equal(t, 1, resources.Remaining(tags.BonusAction), "expected 1 bonus action")
			assert.Equal(t, 1, resources.Remaining(tags.Reaction), "expected 1 reaction")
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
			assert.Equal(t, 0, resources.Remaining(tags.Action), "expected 0 actions")
			assert.False(t, resources.CanUse(tags.Action, 1), "should not be able to use action")
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

			assert.Equal(t, 1, resources.Remaining(tags.Action), "expected 1 action after reset")
			assert.Equal(t, 1, resources.Remaining(tags.BonusAction), "expected 1 bonus action after reset")
			assert.Equal(t, 1, resources.Remaining(tags.Reaction), "expected 1 reaction after reset")
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
			assert.Equal(t, 40, resources.maxSpeed(), "expected max speed 40")

			resources.Consume(tags.WalkSpeed, 10)
			assert.Equal(t, 20, resources.Remaining(tags.WalkSpeed), "expected 20 walk speed remaining")
			assert.Equal(t, 10, resources.Remaining(tags.SwimSpeed), "expected 10 swim speed remaining")
			assert.Equal(t, 30, resources.Remaining(tags.FlySpeed), "expected 30 fly speed remaining")
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
			assert.Equal(t, 5, resources.Remaining(tags.WalkSpeed), "expected 5 walk speed remaining")
			assert.Equal(t, 0, resources.Remaining(tags.SwimSpeed), "expected 0 swim speed remaining")
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

			assert.Equal(t, 30, resources.Remaining(tags.WalkSpeed), "expected 30 walk speed after reset")
			assert.Equal(t, 20, resources.Remaining(tags.SwimSpeed), "expected 20 swim speed after reset")
			assert.Equal(t, 40, resources.Remaining(tags.FlySpeed), "expected 40 fly speed after reset")
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
			assert.Equal(t, 20, resources.Remaining(tags.FlySpeed), "expected 20 fly speed remaining")
			assert.Equal(t, 10, resources.Remaining(tags.WalkSpeed), "expected 10 walk speed remaining")
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

			assert.Equal(t, 3, resources.Remaining(tags.LegendaryAction), "expected 3 legendary actions")
			assert.Equal(t, 5, resources.Remaining(tags.SorceryPoints), "expected 5 sorcery points")

			resources.Consume(tags.LegendaryAction, 2)
			assert.Equal(t, 1, resources.Remaining(tags.LegendaryAction), "expected 1 legendary action after consuming 2")
		})
	})
}
