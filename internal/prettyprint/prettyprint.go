package prettyprint

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

var eventStack []eventbus.Message

func shouldPrintEnd() bool {
	if len(eventStack) == 0 {
		return true
	}

	stoppers := []string{
		core.ExpressionResultType,
		core.CheckResultType,
		core.AttributeCalculationType,
		core.ConfirmType,
		core.AttributeChangedType,
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
	case core.EncounterType:
		return printEncounter(ev.Data.(core.EncounterEvent))
	case core.RoundType:
		return printRound(ev.Data.(core.RoundEvent))
	case core.TurnType:
		return printTurn(ev.Data.(core.TurnEvent))
	case core.DeathType:
		return printDeath(ev.Data.(core.DeathEvent))
	case core.UseActionType:
		return printUseAction(ev.Data.(core.UseActionEvent))
	case core.TakeDamageType:
		return printTakeDamage(ev.Data.(core.TakeDamageEvent))
	case core.ExpressionResultType:
		return printExpressionResult(ev.Data.(core.ExpressionResultEvent))
	case core.CheckResultType:
		return printCheckResult(ev.Data.(core.CheckResultEvent))
	case core.AttackRollType:
		return printAttackRoll(ev.Data.(core.AttackRollEvent))
	case core.AttributeCalculationType:
		return printAttributeCalculation(ev.Data.(core.AttributeCalculationEvent))
	case core.ConfirmType:
		return printConfirm(ev.Data.(core.ConfirmEvent))
	case core.DamageRollType:
		return printDamageRoll(ev.Data.(core.DamageRollEvent))
	case core.EffectType:
		return printEffect(ev.Data.(core.EffectEvent))
	case core.AttributeChangedType:
		return printAttributeChange(ev.Data.(core.AttributeChangeEvent))
	case core.SavingThrowType:
		return printSavingThrow(ev.Data.(core.SavingThrowEvent))
	}
	return "unknown event " + ev.Kind
}

func printWorld(w core.World) string {
	sb := strings.Builder{}
	sb.WriteString("🌍 World\n")
	for y := range w.Height() {
		for x := 0; x < w.Width(); x++ {
			pos := grid.Position{X: x, Y: y}
			nav, _ := w.Navigation.At(pos)
			cell, _ := w.At(pos)
			if !nav.Walkable {
				sb.WriteString("#")
				continue
			}
			if cell.IsOccupied() {
				occupant, _ := cell.Occupant()
				sb.WriteString(occupant.Name[0:1])
				continue
			}
			sb.WriteString(".")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func printActor(a core.Actor) string {
	sb := strings.Builder{}
	stats := []string{
		fmt.Sprintf("HP: %3d/%-3d", a.HitPoints, a.MaxHitPoints),
		fmt.Sprintf("AC: %3d", a.ArmorClass().Value),
	}
	sb.WriteString(fmt.Sprintf("🧝 %-20s %s", a.Name, strings.Join(stats, " ")))
	return sb.String()
}

func printTeam(a []*core.Actor) string {
	sb := strings.Builder{}
	sb.WriteString("🎴 " + a[0].Name)
	out := []string{}
	for _, c := range a {
		out = append(out, indent(printActor(*c)))
	}
	sb.WriteString("\n" + strings.Join(out, "\n"))
	return sb.String()
}

func printEncounter(e core.EncounterEvent) string {
	sb := strings.Builder{}
	sb.WriteString("🏰 Encounter Start")
	sb.WriteString("\n" + indent(printWorld(e.World)))
	teams := map[string][]*core.Actor{}
	for _, c := range e.Actors {
		teams[string(c.Team)] = append(teams[string(c.Team)], c)
	}
	out := []string{}
	for _, t := range teams {
		out = append(out, indent(printTeam(t)))
		out = append(out, "│ └─○")
	}
	sb.WriteString("\n" + strings.Join(out, "\n"))
	sb.WriteString("\n└─○")
	return sb.String()
}

func printRound(r core.RoundEvent) string {
	return fmt.Sprintf("🔄 Round %d", r.Round+1)
}

func printTurn(t core.TurnEvent) string {
	return fmt.Sprintf("🔃 Turn %d: %s", t.Turn+1, t.Actor.Name)
}

func printDeath(d core.DeathEvent) string {
	return fmt.Sprintf("☠️ %s is about to die", d.Actor.Name)
}

func printConfirm(c core.ConfirmEvent) string {
	if c.Confirm {
		return "✅ Confirmed"
	}
	return "❌ Denied"
}

func printUseAction(u core.UseActionEvent) string {
	return fmt.Sprintf("💫 %s uses %s on %s", u.Source.Name, u.Action.Name(), u.Target.Name)
}

func printTakeDamage(d core.TakeDamageEvent) string {
	return fmt.Sprintf("🩸 %s takes %d damage (%d HP left)", d.Target.Name, d.Damage, d.Target.HitPoints)
}

func printExpressionResult(e core.ExpressionResultEvent) string {
	sb := strings.Builder{}
	sb.WriteString("🎲 ")
	sb.WriteString(printExpression(e.Expression))
	return sb.String()
}

func printCheckResult(e core.CheckResultEvent) string {
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

func printAttackRoll(e core.AttackRollEvent) string {
	return fmt.Sprintf("🗡️ %s does an attack roll against %s", e.Source.Name, e.Target.Name)
}

func printAttributeCalculation(e core.AttributeCalculationEvent) string {
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
	sb.WriteString(" ")
	sb.WriteString(printExpression(*e.Expression))
	return sb.String()
}

func printDamageRoll(e core.DamageRollEvent) string {
	return fmt.Sprintf("🎲 %s is rolling damage", e.Source.Name)
}

func printEffect(e core.EffectEvent) string {
	return fmt.Sprintf("⚡ %s triggered", e.Effect.Name)
}

func printAttributeChange(e core.AttributeChangeEvent) string {
	return fmt.Sprintf("🔀 %s %s changed from %d to %d", e.Source.Name, tags.ToReadable(e.Attribute), e.OldValue, e.Value)
}

func printSavingThrow(e core.SavingThrowEvent) string {
	return fmt.Sprintf("🍥 %s rolls a %s saving throw against DC %d", e.Source.Name, tags.ToReadable(e.Attribute), e.DifficultyClass)
}
