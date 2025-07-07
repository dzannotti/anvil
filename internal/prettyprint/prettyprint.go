package prettyprint

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"anvil/internal/core"
	"anvil/internal/eventbus"
)

type EventFormatter func(eventbus.Message) string

var eventStack []eventbus.Message

func shouldPrintEnd() bool {
	if len(eventStack) == 0 {
		return true
	}

	stoppers := []string{
		core.ExpressionResultType,
		core.CheckResultType,
		core.SavingThrowResultType,
		core.AttributeCalculationType,
		core.ConfirmType,
		core.AttributeChangedType,
		core.SpendResourceType,
		core.ConditionChangedType,
		core.DeathSavingThrowResultType,
		core.TargetType,
	}

	lastEvent := eventStack[len(eventStack)-1]
	return !slices.Contains(stoppers, lastEvent.Kind)
}

func Print(out io.Writer, event eventbus.Message) {
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
	return func(event eventbus.Message) string {
		data := event.Data.(T)
		return typedFormatter(data)
	}
}
