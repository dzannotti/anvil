package prettyprint

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/tag"
	"fmt"
	"slices"
	"strings"
)

var eventFormatters = map[string]EventFormatter{
	core.EncounterType:                 makeFormatter(printEncounter),
	core.RoundType:                     makeFormatter(printRound),
	core.TurnType:                      makeFormatter(printTurn),
	core.DeathType:                     makeFormatter(printDeath),
	core.UseActionType:                 makeFormatter(printUseAction),
	core.TakeDamageType:                makeFormatter(printTakeDamage),
	core.ExpressionResultType:          makeFormatter(printExpressionResult),
	core.CheckResultType:               makeFormatter(printCheckResult),
	core.AttackRollType:                makeFormatter(printAttackRoll),
	core.AttributeCalculationType:      makeFormatter(printAttributeCalculation),
	core.ConfirmType:                   makeFormatter(printConfirm),
	core.DamageRollType:                makeFormatter(printDamageRoll),
	core.EffectType:                    makeFormatter(printEffect),
	core.AttributeChangedType:          makeFormatter(printAttributeChange),
	core.SavingThrowType:               makeFormatter(printSavingThrow),
	core.SpendResourceType:             makeFormatter(printSpendResource),
	core.ConditionChangedType:          makeFormatter(printConditionChanged),
	core.MoveType:                      makeFormatter(printMove),
	core.MoveStepType:                  makeFormatter(printMoveStep),
	core.DeathSavingThrowType:          makeFormatter(printDeathSavingThrow),
	core.DeathSavingThrowResultType:    makeFormatter(printDeathSavingThrowResult),
	core.DeathSavingThrowAutomaticType: makeFormatter(printDeathSavingThrowAutomaticResult),
	core.SavingThrowResultType:         makeFormatter(printSavingThrowResult),
	core.TargetType:                    makeFormatter(printTarget),
}

func formatEvent(event eventbus.Message) string {
	eventType := event.Kind
	if formatter, exists := eventFormatters[eventType]; exists {
		return formatter(event)
	}

	return fmt.Sprintf("📝 %s: %v", eventType, event.Data)
}

func printPosition(pos grid.Position) string {
	return fmt.Sprintf("(%d, %d)", pos.X, pos.Y)
}

func printWorld(w *core.World, path []grid.Position) string {
	sb := strings.Builder{}
	sb.WriteString("🌍 World\n")
	for y := range w.Height() {
		for x := range w.Width() {
			pos := grid.Position{X: x, Y: y}
			cell, _ := w.At(pos)
			if cell.Tile == core.Wall {
				sb.WriteString("#")
				continue
			}
			if cell.IsOccupied() {
				occupant, _ := cell.Occupant()
				sb.WriteString(occupant.Name[0:1])
				continue
			}
			if slices.Contains(path, pos) {
				sb.WriteString("*")
				continue
			}
			sb.WriteString(".")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func printActor(a *core.Actor) string {
	sb := strings.Builder{}
	stats := []string{
		fmt.Sprintf("HP: %3d/%-3d", a.HitPoints, a.MaxHitPoints),
		fmt.Sprintf("AC: %3d", a.ArmorClass().Value),
	}
	sb.WriteString(fmt.Sprintf("🧝 %-20s %s", a.Name, strings.Join(stats, " ")))
	return sb.String()
}

func printActors(a []*core.Actor) string {
	sb := strings.Builder{}
	out := []string{}
	for _, c := range a {
		out = append(out, printActor(c))
	}
	sb.WriteString(strings.Join(out, "\n"))
	return sb.String()
}

func printTeam(a []*core.Actor) string {
	sb := strings.Builder{}
	sb.WriteString("🎴 " + string(a[0].Team))
	sb.WriteString("\n")
	sb.WriteString(indent(printActors(a), 0))
	return sb.String()
}

func printEncounter(e core.EncounterEvent) string {
	sb := strings.Builder{}
	sb.WriteString("🏰 Encounter Start")
	sb.WriteString("\n" + indent(printWorld(e.World, []grid.Position{}), 0))
	teams := map[string][]*core.Actor{}
	for _, c := range e.Actors {
		teams[string(c.Team)] = append(teams[string(c.Team)], c)
	}
	out := []string{}
	for _, t := range teams {
		out = append(out, indent(printTeam(t), 1))
		out = append(out, "│ └─○")
	}
	sb.WriteString("\n" + strings.Join(out, "\n"))
	sb.WriteString("\n└─○")
	return sb.String()
}

func printRound(r core.RoundEvent) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("🔄 Round %d", r.Round+1))
	sb.WriteString("\n")
	sb.WriteString(indent(printActors(r.Actors), 1))
	return sb.String()
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
	pos := make([]string, len(u.Target))
	for i, p := range u.Target {
		pos[i] = printPosition(p)
	}
	return fmt.Sprintf("💫 %s uses %s at [%s]", u.Source.Name, u.Action.Name(), strings.Join(pos, ", "))
}

func printTakeDamage(d core.TakeDamageEvent) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("🩸 %s takes %d damage (%d HP left)", d.Target.Name, d.Damage.Value, d.Target.HitPoints))
	sb.WriteString("\n")
	sb.WriteString(indent(printExpression(d.Damage, true), 1))
	return sb.String()
}

func printExpressionResult(e core.ExpressionResultEvent) string {
	sb := strings.Builder{}
	sb.WriteString("🎲 ")
	sb.WriteString(printExpression(e.Expression))
	return sb.String()
}

func printSavingThrowResult(e core.SavingThrowResultEvent) string {
	return formatRollResult(e.Success, e.Critical, e.Value, e.Against)
}

func formatRollResult(success, critical bool, value, against int) string {
	icons := map[bool]string{true: "✅", false: "❌"}
	sb := strings.Builder{}
	sb.WriteString(icons[success])
	if critical {
		sb.WriteString("💥 Critical")
	}
	sb.WriteString(map[bool]string{true: " Success", false: " Failure"}[success])
	return fmt.Sprintf("%s %d vs %d", sb.String(), value, against)
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
	sb.WriteString(printExpression(e.Expression))
	return sb.String()
}

func printDamageRoll(e core.DamageRollEvent) string {
	return fmt.Sprintf("🎲 %s is rolling damage", e.Source.Name)
}

func printEffect(e core.EffectEvent) string {
	return fmt.Sprintf("⚡ %s triggered", e.Effect.Name)
}

func printAttributeChange(e core.AttributeChangeEvent) string {
	return fmt.Sprintf(
		"🔀 %s %s changed from %d to %d",
		e.Source.Name,
		tags.ToReadable(e.Attribute),
		e.OldValue,
		e.Value,
	)
}

func printSavingThrow(e core.SavingThrowEvent) string {
	return fmt.Sprintf(
		"🍥 %s rolls a %s saving throw against DC %d",
		e.Source.Name,
		tags.ToReadable(e.Attribute),
		e.DifficultyClass,
	)
}

func printSpendResource(e core.SpendResourceEvent) string {
	return fmt.Sprintf("🧾 %s spent %d %s", e.Source.Name, e.Amount, tags.ToReadable(e.Resource))
}

func printConditionChanged(e core.ConditionChangedEvent) string {
	emoji := "➕"
	text := "gains condition"
	if !e.Added {
		emoji = "➖"
		text = "loses condition"
	}
	from := ""
	if e.From != nil {
		from = fmt.Sprintf("from %s", e.From.Name)
	}
	return fmt.Sprintf("%s %s %s %s %s", emoji, e.Source.Name, text, tags.ToReadable(e.Condition), from)
}

func printMove(e core.MoveEvent) string {
	sb := strings.Builder{}
	sb.WriteString(
		fmt.Sprintf("🚶 %s wants to move from %s to %s", e.Source.Name, printPosition(e.From), printPosition(e.To)),
	)
	sb.WriteString("\n")
	sb.WriteString(indent(printWorld(e.World, e.Path.Path), 1))
	return sb.String()
}

func printMoveStep(e core.MoveStepEvent) string {
	return fmt.Sprintf("🚶 %s about to step from %s to %s", e.Source.Name, printPosition(e.From), printPosition(e.To))
}

func printDeathSavingThrow(e core.DeathSavingThrowEvent) string {
	return fmt.Sprintf("⚰️ %s is about to roll a Death Saving throw", e.Source.Name)
}

func printDeathSavingThrowResult(e core.DeathSavingThrowResultEvent) string {
	return fmt.Sprintf("%d failures and %d successes", e.Failure, e.Success)
}

func printDeathSavingThrowAutomaticResult(e core.DeathSavingThrowAutomaticEvent) string {
	status := "success"
	if e.Failure {
		status = "failure"
	}
	return fmt.Sprintf("⚰️ %s automatic death save throw %s", e.Source.Name, status)
}

func printTarget(e core.TargetEvent) string {
	targets := make([]string, len(e.Target))
	for i, t := range e.Target {
		targets[i] = t.Name
	}
	return fmt.Sprintf("🎯 Target [%s]", strings.Join(targets, ", "))
}
