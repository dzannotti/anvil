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
					tags.ResourceAction:      1,
					tags.ResourceBonusAction: 1,
					tags.ResourceReaction:    1,
				},
			}
			resources.Reset()
			resources.LongRest()
			assert.Equal(t, 1, resources.Remaining(tags.ResourceAction), "expected 1 action")
			assert.Equal(t, 1, resources.Remaining(tags.ResourceBonusAction), "expected 1 bonus action")
			assert.Equal(t, 1, resources.Remaining(tags.ResourceReaction), "expected 1 reaction")
		})

		t.Run("should consume actions when used", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceAction:      1,
					tags.ResourceBonusAction: 1,
					tags.ResourceReaction:    1,
				},
			}
			resources.Reset()
			resources.LongRest()
			resources.Consume(tags.ResourceAction, 1)
			assert.Equal(t, 0, resources.Remaining(tags.ResourceAction), "expected 0 actions")
			assert.False(t, resources.CanUse(tags.ResourceAction, 1), "should not be able to use action")
		})

		t.Run("should reset all actions at the start of new turn", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceAction:      1,
					tags.ResourceBonusAction: 1,
					tags.ResourceReaction:    1,
				},
			}
			resources.Reset()
			resources.LongRest()
			resources.Consume(tags.ResourceAction, 1)
			resources.Consume(tags.ResourceBonusAction, 1)
			resources.Consume(tags.ResourceReaction, 1)

			resources.Reset()

			assert.Equal(t, 1, resources.Remaining(tags.ResourceAction), "expected 1 action after reset")
			assert.Equal(t, 1, resources.Remaining(tags.ResourceBonusAction), "expected 1 bonus action after reset")
			assert.Equal(t, 1, resources.Remaining(tags.ResourceReaction), "expected 1 reaction after reset")
		})
	})

	t.Run("Movement Speed", func(t *testing.T) {
		t.Run("should calculate max speed and track remaining movement correctly", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceWalkSpeed: 30,
					tags.ResourceSwimSpeed: 20,
					tags.ResourceFlySpeed:  40,
				},
			}
			resources.Reset()
			resources.LongRest()
			assert.Equal(t, 40, resources.maxSpeed(), "expected max speed 40")

			resources.Consume(tags.ResourceWalkSpeed, 10)
			assert.Equal(t, 20, resources.Remaining(tags.ResourceWalkSpeed), "expected 20 walk speed remaining")
			assert.Equal(t, 10, resources.Remaining(tags.ResourceSwimSpeed), "expected 10 swim speed remaining")
			assert.Equal(t, 30, resources.Remaining(tags.ResourceFlySpeed), "expected 30 fly speed remaining")
		})

		t.Run("should not allow movement beyond available speed", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceWalkSpeed: 30,
					tags.ResourceSwimSpeed: 20,
					tags.ResourceFlySpeed:  40,
				},
			}
			resources.Reset()

			resources.Consume(tags.ResourceWalkSpeed, 25)
			assert.Equal(t, 5, resources.Remaining(tags.ResourceWalkSpeed), "expected 5 walk speed remaining")
			assert.Equal(t, 0, resources.Remaining(tags.ResourceSwimSpeed), "expected 0 swim speed remaining")
		})

		t.Run("should reset movement speed on new turn", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceWalkSpeed: 30,
					tags.ResourceSwimSpeed: 20,
					tags.ResourceFlySpeed:  40,
				},
			}
			resources.Reset()

			resources.Consume(tags.ResourceWalkSpeed, 15)
			resources.Reset()

			assert.Equal(t, 30, resources.Remaining(tags.ResourceWalkSpeed), "expected 30 walk speed after reset")
			assert.Equal(t, 20, resources.Remaining(tags.ResourceSwimSpeed), "expected 20 swim speed after reset")
			assert.Equal(t, 40, resources.Remaining(tags.ResourceFlySpeed), "expected 40 fly speed after reset")
		})

		t.Run("should handle different movement types independently", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceWalkSpeed: 30,
					tags.ResourceSwimSpeed: 20,
					tags.ResourceFlySpeed:  40,
				},
			}
			resources.Reset()

			resources.Consume(tags.ResourceFlySpeed, 20)
			assert.Equal(t, 20, resources.Remaining(tags.ResourceFlySpeed), "expected 20 fly speed remaining")
			assert.Equal(t, 10, resources.Remaining(tags.ResourceWalkSpeed), "expected 10 walk speed remaining")
		})
	})

	t.Run("Custom Resources", func(t *testing.T) {
		t.Run("should handle custom resources", func(t *testing.T) {
			resources := Resources{
				Max: map[tag.Tag]int{
					tags.ResourceLegendaryAction: 3,
					tags.ResourceSorceryPoints:   5,
				},
			}

			resources.Reset()
			resources.LongRest()

			assert.Equal(t, 3, resources.Remaining(tags.ResourceLegendaryAction), "expected 3 legendary actions")
			assert.Equal(t, 5, resources.Remaining(tags.ResourceSorceryPoints), "expected 5 sorcery points")

			resources.Consume(tags.ResourceLegendaryAction, 2)
			assert.Equal(t, 1, resources.Remaining(tags.ResourceLegendaryAction), "expected 1 legendary action after consuming 2")
		})
	})
}
