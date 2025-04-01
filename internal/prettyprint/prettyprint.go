package prettyprint

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"anvil/internal/core"
	"anvil/internal/core/definition"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
)

var eventStack []eventbus.Message

func shouldPrintEnd() bool {
	if len(eventStack) == 0 {
		return true
	}

	stoppers := []string{
		core.ExpressionResultEventType,
		core.CheckResultEventType,
		core.AttributeCalculationEventType,
	}

	lastEvent := eventStack[len(eventStack)-1]
	return !slices.Contains(stoppers, lastEvent.Kind)
}

func Print(out io.Writer, ev eventbus.Message) {
	depthPrefix := strings.Repeat("│  ", max(0, ev.Depth-1))
	if ev.End {
		if shouldPrintEnd() {
			fmt.Fprintln(out, depthPrefix+"└─○")
		}
		eventStack = eventStack[:len(eventStack)-1]
		return
	}
	if !ev.End {
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
	switch ev.Kind {
	case core.EncounterEventType:
		return printEncounter(ev.Data.(core.EncounterEvent))
		/*case Round:
			return printRound(e)
		case Turn:
			return printTurn(e)
		case Died:
			return printDeath(e)
		case UseAction:
			return printUseAction(e)
		case TakeDamage:
			return printTakeDamage(e)
		case ExpressionResult:
			return printExpressionResult(e)
		case CheckResult:
			return printCheckResult(e)
		case AttackRoll:
			return printAttackRoll(e)
		case AttributeCalculation:
			return printAttributeCalculation(e)
		}*/
	}
	return "unknown event " + ev.Kind
}

func printWorld(w definition.World) string {
	sb := strings.Builder{}
	sb.WriteString("🌍 World\n")
	for x := range w.Width() {
		for y := range w.Height() {
			pos := grid.Position{X: x, Y: y}
			nav, _ := w.Navigation().At(pos)
			cell, _ := w.At(pos)
			if !nav.IsWalkable() {
				sb.WriteString("#")
				continue
			}
			if cell.IsOccupied() {
				occupant, _ := cell.Occupant()
				sb.WriteString(occupant.Name()[0:1])
				continue
			}
			sb.WriteString(".")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func printCreature(c definition.Creature) string {
	sb := strings.Builder{}
	stats := []string{
		fmt.Sprintf("HP: %3d/%d", c.HitPoints(), c.MaxHitPoints()),
	}
	sb.WriteString(fmt.Sprintf("🧝 %-20s %s", c.Name(), strings.Join(stats, " ")))
	return sb.String()
}

func printTeam(t definition.Team) string {
	sb := strings.Builder{}
	sb.WriteString("🎴 " + t.Name())
	creatures := []string{}
	for _, c := range t.Members() {
		creatures = append(creatures, indent(printCreature(c)))
	}
	sb.WriteString("\n" + strings.Join(creatures, "\n"))
	return sb.String()
}

func printEncounter(e core.EncounterEvent) string {
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

/*
func printRound(r RoundEvent) string {
	return fmt.Sprintf("🔄 Round %d", r.Round+1)
}

func printTurn(t TurnEvent) string {
	return fmt.Sprintf("🔃 Turn %d: %s", t.Turn+1, t.Creature.Name)
}

func printDeath(d Died) string {
	return fmt.Sprintf("☠️ %s is about to die", d.Creature.Name)
}

func printUseAction(u UseAction) string {
	return fmt.Sprintf("💫 %s uses %s on %s", u.Source.Name, u.Action.Name, u.Target.Name)
}

func printTakeDamage(d TakeDamage) string {
	return fmt.Sprintf("🩸 %s takes %d damage", d.Target.Name, d.Damage)
}

func printExpressionResult(e ExpressionResult) string {
	sb := strings.Builder{}
	sb.WriteString("🎲 ")
	sb.WriteString(printExpression(e.Expression))
	return sb.String()
}

func printCheckResult(e CheckResult) string {
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

func printAttackRoll(e AttackRoll) string {
	return fmt.Sprintf("🗡️ %s does an attack roll against %s", e.Source.Name, e.Target.Name)
}

func printAttributeCalculation(e AttributeCalculation) string {
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
*/
