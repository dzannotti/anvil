package prettyprint

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"anvil/internal/core/event"
	"anvil/internal/core/event/snapshot"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/tag"
)

var eventStack []eventbus.Message

func shouldPrintEnd() bool {
	if len(eventStack) == 0 {
		return true
	}

	stoppers := []reflect.Type{
		reflect.TypeOf(event.ExpressionResult{}),
		reflect.TypeOf(event.CheckResult{}),
		reflect.TypeOf(event.AttributeCalculation{}),
		/*reflect.TypeOf(event.AttributeChange{}),
		reflect.TypeOf(event.Target{}),
		reflect.TypeOf(event.SpendResource{}),*/
	}

	lastEvent := eventStack[len(eventStack)-1]
	lastEventType := reflect.TypeOf(lastEvent.Data)

	for _, stopper := range stoppers {
		if lastEventType == stopper {
			return false
		}
	}
	return true
}

func Print(out io.Writer, ev eventbus.Message) {
	depthPrefix := strings.Repeat("│  ", max(0, ev.Depth-1))
	if ev.IsEnd {
		if shouldPrintEnd() {
			fmt.Fprintln(out, depthPrefix+"└─○")
		}
		eventStack = eventStack[:len(eventStack)-1]
		return
	}
	if !ev.IsEnd {
		eventStack = append(eventStack, ev)
	}
	extraPrefix := ""
	if ev.Depth > 0 {
		extraPrefix = "├─ "
	}
	eventString := printMessage(ev)
	lines := strings.Split(eventString, "\n")
	first := depthPrefix + extraPrefix + lines[0]
	fmt.Fprintln(out, first)
	for _, line := range lines[1:] {
		fmt.Fprintln(out, depthPrefix+"│  "+line)
	}
}

func printMessage(ev eventbus.Message) string {
	switch e := ev.Data.(type) {
	case event.Encounter:
		return printEncounter(e)
	case event.Round:
		return printRound(e)
	case event.Turn:
		return printTurn(e)
	case event.Died:
		return printDeath(e)
	case event.UseAction:
		return printUseAction(e)
	case event.TakeDamage:
		return printTakeDamage(e)
	case event.ExpressionResult:
		return printExpressionResult(e)
	case event.CheckResult:
		return printCheckResult(e)
	case event.AttackRoll:
		return printAttackRoll(e)
	case event.AttributeCalculation:
		return printAttributeCalculation(e)
	}
	t := reflect.TypeOf(ev.Data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return "unknown event " + t.Name()
}

func printWorld(w snapshot.World) string {
	sb := strings.Builder{}
	sb.WriteString("🌍 World\n")
	for x := range len(w.Cells) {
		for y := range len(w.Cells[x]) {
			if !w.Cells[x][y].Walkable {
				sb.WriteString("#")
				continue
			}
			if w.Cells[x][y].Occupant.Name != "" {
				sb.WriteString(w.Cells[x][y].Occupant.Name[0:1])
				continue
			}
			sb.WriteString(".")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func printCreature(c snapshot.Creature) string {
	sb := strings.Builder{}
	stats := []string{
		fmt.Sprintf("HP: %3d/%3d", c.HitPoints, c.MaxHitPoints),
	}
	sb.WriteString(fmt.Sprintf("🧝 %-20s %s", c.Name, strings.Join(stats, " ")))
	return sb.String()
}

func printTeam(f snapshot.Team) string {
	sb := strings.Builder{}
	sb.WriteString("🎴 " + f.Name)
	creatures := []string{}
	for _, c := range f.Members {
		creatures = append(creatures, indent(printCreature(c)))
	}
	sb.WriteString("\n" + strings.Join(creatures, "\n"))
	return sb.String()
}

func printEncounter(e event.Encounter) string {
	sb := strings.Builder{}
	sb.WriteString("🏰 Encounter Start")
	teams := []string{}
	sb.WriteString("\n" + indent(printWorld(e.World)))
	for _, f := range e.Teams {
		teams = append(teams, indent(printTeam(f)))
		teams = append(teams, "│ └─○")
	}
	sb.WriteString("\n" + strings.Join(teams, "\n"))
	sb.WriteString("\n└─○")
	return sb.String()
}

func printRound(r event.Round) string {
	return fmt.Sprintf("🔄 Round %d", r.Round+1)
}

func printTurn(t event.Turn) string {
	return fmt.Sprintf("🔃 Turn %d: %s", t.Turn+1, t.Creature.Name)
}

func printDeath(d event.Died) string {
	return fmt.Sprintf("☠️ %s is about to die", d.Creature.Name)
}

func printUseAction(u event.UseAction) string {
	return fmt.Sprintf("💫 %s uses %s on %s", u.Source.Name, u.Action.Name, u.Target.Name)
}

func printTakeDamage(d event.TakeDamage) string {
	return fmt.Sprintf("🩸 %s takes %d damage", d.Target.Name, d.Damage)
}

func printExpressionResult(e event.ExpressionResult) string {
	sb := strings.Builder{}
	sb.WriteString("🎲 ")
	sb.WriteString(printExpression(e.Expression))
	return sb.String()
}

func printCheckResult(e event.CheckResult) string {
	sb := strings.Builder{}
	sIcon := map[bool]string{true: "✅", false: "❌"}
	sb.WriteString(sIcon[e.Success])
	if e.Critical {
		sb.WriteString("💥 Critical")
	}
	success := map[bool]string{true: " Success", false: " Failure"}
	sb.WriteString(success[e.Success])
	outcome := sb.String()
	return fmt.Sprintf("%s %d vs %d", outcome, e.Value, e.Against)
}

func printAttackRoll(e event.AttackRoll) string {
	return fmt.Sprintf("🗡️ %s does an attack roll against %s", e.Source.Name, e.Target.Name)
}

func printAttributeCalculation(e event.AttributeCalculation) string {
	emoji := map[tag.Tag]string{
		tags.ArmorClass:   "🛡️",
		tags.Strength:     "💪",
		tags.Dexterity:    "🏹",
		tags.Constitution: "❤️",
		tags.Intelligence: "🧠",
		tags.Wisdom:       "🧘",
		tags.Charisma:     "👑",
	}
	sb := strings.Builder{}
	sb.WriteString(emoji[e.Attribute])
	sb.WriteString("\n")
	sb.WriteString(printExpression(e.Expression))
	return sb.String()
}
