package prettyprint

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/tag"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

var eventFormatters = map[string]EventFormatter{
	reflect.TypeOf(core.EncounterEvent{}).String():                 makeFormatter(printEncounter),
	reflect.TypeOf(core.RoundEvent{}).String():                     makeFormatter(printRound),
	reflect.TypeOf(core.TurnEvent{}).String():                      makeFormatter(printTurn),
	reflect.TypeOf(core.DeathEvent{}).String():                     makeFormatter(printDeath),
	reflect.TypeOf(core.UseActionEvent{}).String():                 makeFormatter(printUseAction),
	reflect.TypeOf(core.TakeDamageEvent{}).String():                makeFormatter(printTakeDamage),
	reflect.TypeOf(core.ExpressionResultEvent{}).String():          makeFormatter(printExpressionResult),
	reflect.TypeOf(core.CheckResultEvent{}).String():               makeFormatter(printCheckResult),
	reflect.TypeOf(core.AttackRollEvent{}).String():                makeFormatter(printAttackRoll),
	reflect.TypeOf(core.AttributeCalculationEvent{}).String():      makeFormatter(printAttributeCalculation),
	reflect.TypeOf(core.ConfirmEvent{}).String():                   makeFormatter(printConfirm),
	reflect.TypeOf(core.DamageRollEvent{}).String():                makeFormatter(printDamageRoll),
	reflect.TypeOf(core.EffectEvent{}).String():                    makeFormatter(printEffect),
	reflect.TypeOf(core.AttributeChangeEvent{}).String():           makeFormatter(printAttributeChange),
	reflect.TypeOf(core.SavingThrowEvent{}).String():               makeFormatter(printSavingThrow),
	reflect.TypeOf(core.SpendResourceEvent{}).String():             makeFormatter(printSpendResource),
	reflect.TypeOf(core.ConditionChangedEvent{}).String():          makeFormatter(printConditionChanged),
	reflect.TypeOf(core.MoveEvent{}).String():                      makeFormatter(printMove),
	reflect.TypeOf(core.MoveStepEvent{}).String():                  makeFormatter(printMoveStep),
	reflect.TypeOf(core.DeathSavingThrowEvent{}).String():          makeFormatter(printDeathSavingThrow),
	reflect.TypeOf(core.DeathSavingThrowResultEvent{}).String():    makeFormatter(printDeathSavingThrowResult),
	reflect.TypeOf(core.DeathSavingThrowAutomaticEvent{}).String(): makeFormatter(printDeathSavingThrowAutomaticResult),
	reflect.TypeOf(core.SavingThrowResultEvent{}).String():         makeFormatter(printSavingThrowResult),
	reflect.TypeOf(core.TargetEvent{}).String():                    makeFormatter(printTarget),
}

func formatEvent(event eventbus.Event) string {
	eventType := event.Kind
	if formatter, exists := eventFormatters[eventType]; exists {
		return formatter(event)
	}

	return fmt.Sprintf("📝 %s: %v", eventType, event.Data)
}

func printPosition(pos grid.Position) string {
	return fmt.Sprintf("(%d, %d)", pos.X, pos.Y)
}

func printPositions(positions []grid.Position) string {
	pos := make([]string, len(positions))
	for i, p := range positions {
		pos[i] = printPosition(p)
	}
	return strings.Join(pos, ", ")
}

func printActorNames(actors []*core.Actor) string {
	names := make([]string, len(actors))
	for i, a := range actors {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

func printWorld(w *core.World, path []grid.Position) string {
	sb := strings.Builder{}
	sb.WriteString("🌍 World\n")
	for y := range w.Height() {
		for x := range w.Width() {
			pos := grid.Position{X: x, Y: y}
			cell := w.At(pos)
			if cell == nil {
				sb.WriteString(".")
				continue
			}

			if cell.Tile == core.Wall {
				sb.WriteString("#")
				continue
			}

			if cell.IsOccupied() {
				occupant := cell.Occupant()
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
	tb := NewTreeBuilder()
	tb.AddRawLine("🏰 Encounter Start")

	// Add world section - no closure, world continues
	tb.AddIndentedBlock(printWorld(e.World, []grid.Position{}))

	// Add teams
	teams := map[string][]*core.Actor{}
	for _, c := range e.Actors {
		teams[string(c.Team)] = append(teams[string(c.Team)], c)
	}

	for _, t := range teams {
		tb.AddIndentedBlock(printTeam(t))
		// Team closure should be indented under the team
		tb.AddRawLine(TreeVertical + TreeEndCircle)
	}

	// No final closure - encounter stays open for rounds
	return tb.String()
}

func printRound(r core.RoundEvent) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("🔄 Round %d", r.Round+1))
	sb.WriteString("\n")
	sb.WriteString(indent(printActors(r.Actors), 0))
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
	return fmt.Sprintf("💫 %s uses %s at [%s]", u.Source.Name, u.Action.Name(), printPositions(u.Target))
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
	return formatRollResult(e.Success, e.Critical, e.Value, e.Against)
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
	return fmt.Sprintf("🎯 Target [%s]", printActorNames(e.Target))
}
