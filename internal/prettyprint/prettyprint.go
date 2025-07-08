package prettyprint

import (
	"fmt"
	"io"
	"reflect"
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
		reflect.TypeOf(core.ExpressionResultEvent{}).String(),
		reflect.TypeOf(core.CheckResultEvent{}).String(),
		reflect.TypeOf(core.SavingThrowResultEvent{}).String(),
		reflect.TypeOf(core.AttributeCalculationEvent{}).String(),
		reflect.TypeOf(core.ConfirmEvent{}).String(),
		reflect.TypeOf(core.AttributeChangeEvent{}).String(),
		reflect.TypeOf(core.SpendResourceEvent{}).String(),
		reflect.TypeOf(core.ConditionChangedEvent{}).String(),
		reflect.TypeOf(core.DeathSavingThrowResultEvent{}).String(),
		reflect.TypeOf(core.TargetEvent{}).String(),
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
