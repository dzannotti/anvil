package prettyprint

import (
	"testing"

	"anvil/internal/core"
	"anvil/internal/eventbus"
	"anvil/internal/tag"

	"github.com/stretchr/testify/assert"
)

func TestEventFormatters_NoErrors(t *testing.T) {
	// Create sample data for testing
	sampleActor := &core.Actor{
		Name:         "Test Actor",
		HitPoints:    10,
		MaxHitPoints: 20,
	}

	tests := []struct {
		name      string
		eventType string
		eventData interface{}
	}{
		{
			name:      "confirm event",
			eventType: eventbus.EventType(core.ConfirmEvent{}),
			eventData: core.ConfirmEvent{Confirm: true},
		},
		{
			name:      "target event",
			eventType: eventbus.EventType(core.TargetEvent{}),
			eventData: core.TargetEvent{Target: []*core.Actor{sampleActor}},
		},
		{
			name:      "check result event",
			eventType: eventbus.EventType(core.CheckResultEvent{}),
			eventData: core.CheckResultEvent{
				Success:  true,
				Critical: false,
				Value:    15,
				Against:  10,
			},
		},
		{
			name:      "saving throw result event",
			eventType: eventbus.EventType(core.SavingThrowResultEvent{}),
			eventData: core.SavingThrowResultEvent{
				Success:  false,
				Critical: true,
				Value:    8,
				Against:  12,
			},
		},
		{
			name:      "attack roll event",
			eventType: eventbus.EventType(core.AttackRollEvent{}),
			eventData: core.AttackRollEvent{
				Source: sampleActor,
				Target: sampleActor,
			},
		},
		{
			name:      "damage roll event",
			eventType: eventbus.EventType(core.DamageRollEvent{}),
			eventData: core.DamageRollEvent{Source: sampleActor},
		},
		{
			name:      "effect event",
			eventType: eventbus.EventType(core.EffectEvent{}),
			eventData: core.EffectEvent{Effect: &core.Effect{Name: "Test Effect"}},
		},
		{
			name:      "spend resource event",
			eventType: eventbus.EventType(core.SpendResourceEvent{}),
			eventData: core.SpendResourceEvent{
				Source:   sampleActor,
				Amount:   1,
				Resource: tag.Tag{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := eventbus.Event{
				Kind: tt.eventType,
				Data: tt.eventData,
			}

			// Should not panic and should return a non-empty string
			assert.NotPanics(t, func() {
				result := formatEvent(msg)
				assert.NotEmpty(t, result, "Event formatter should return non-empty result")
			})
		})
	}
}

func TestFormatEvent_UnknownType(t *testing.T) {
	msg := eventbus.Event{
		Kind: "unknown_event_type",
		Data: "some data",
	}

	result := formatEvent(msg)
	assert.Contains(t, result, "📝 unknown_event_type")
	assert.Contains(t, result, "some data")
}

func TestShouldPrintEnd(t *testing.T) {
	// Clear event stack
	eventStack = []eventbus.Event{}

	// Empty stack should print end
	assert.True(t, shouldPrintEnd())

	// Add a regular event
	eventStack = append(eventStack, eventbus.Event{Kind: eventbus.EventType(core.TurnEvent{})})
	assert.True(t, shouldPrintEnd())

	// Add a stopper event
	eventStack = append(eventStack, eventbus.Event{Kind: eventbus.EventType(core.ConfirmEvent{})})
	assert.False(t, shouldPrintEnd())

	// Clear for other tests
	eventStack = []eventbus.Event{}
}
