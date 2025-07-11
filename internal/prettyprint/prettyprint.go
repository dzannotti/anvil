package prettyprint

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"anvil/internal/core"
	"anvil/internal/eventbus"
)

type EventFormatter func(eventbus.Event) string

var eventStack []eventbus.Event

func shouldPrintEnd() bool {
	if len(eventStack) == 0 {
		return true
	}

	stoppers := []string{
		eventbus.EventType(core.ExpressionResultEvent{}),
		eventbus.EventType(core.CheckResultEvent{}),
		eventbus.EventType(core.SavingThrowResultEvent{}),
		eventbus.EventType(core.AttributeCalculationEvent{}),
		eventbus.EventType(core.ConfirmEvent{}),
		eventbus.EventType(core.AttributeChangeEvent{}),
		eventbus.EventType(core.SpendResourceEvent{}),
		eventbus.EventType(core.ConditionChangedEvent{}),
		eventbus.EventType(core.DeathSavingThrowResultEvent{}),
		eventbus.EventType(core.TargetEvent{}),
	}

	lastEvent := eventStack[len(eventStack)-1]
	return !slices.Contains(stoppers, lastEvent.Kind)
}

func Print(out io.Writer, event eventbus.Event) {
	if event.End {
		if shouldPrintEnd() {
			depth := strings.Repeat(TreeVertical, max(0, event.Depth))
			fmt.Fprintf(out, "%s%s\n", depth, TreeEndCircle)
		}

		if len(eventStack) > 0 {
			eventStack = eventStack[:len(eventStack)-1]
		}

		return
	}

	eventStack = append(eventStack, event)
	eventString := formatEvent(event)
	if eventString == "" {
		return
	}

	if event.Depth == 0 {
		fmt.Fprintf(out, "%s\n", eventString)
		return
	}

	fmt.Fprintf(out, "%s\n", indent(eventString, event.Depth-1))
}

func makeFormatter[T any](typedFormatter func(T) string) EventFormatter {
	return func(event eventbus.Event) string {
		data := event.Data.(T)
		return typedFormatter(data)
	}
}
