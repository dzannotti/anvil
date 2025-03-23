package prettyprint

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"anvil/internal/core/event"
	"anvil/internal/core/team"
	"anvil/internal/log"
)

func Print(out io.Writer, ev log.Event) {
	depthPrefix := strings.Repeat("│  ", max(0, ev.Depth-1))
	if ev.IsEnd {
		fmt.Fprintln(out, depthPrefix+"└─○")
		return
	}
	extraPrefix := ""
	if ev.Depth > 0 {
		extraPrefix = "├─ "
	}
	fmt.Fprintln(out, depthPrefix+extraPrefix+eventToString(ev))
}

func eventToString(ev log.Event) string {
	switch e := ev.Data.(type) {
	case *event.Encounter:
		return printEncounter(*e)
	case *event.Round:
		return printRound(*e)
	case *event.Turn:
		return printTurn(*e)
	case *event.Death:
		return printDeath(*e)
	case *event.UseAction:
		return printUseAction(*e)
	case *event.TakeDamage:
		return printTakeDamage(*e)
	}
	return "unknown event" + reflect.TypeOf(ev.Data).Name()
}

func printFaction(f team.Team, c []*event.Creature) string {
	sb := strings.Builder{}
	sb.WriteString("🎴 " + f.String())
	sb.WriteString(fmt.Sprint(f))
	return sb.String()
}

func printEncounter(e event.Encounter) string {
	sb := strings.Builder{}
	sb.WriteString("🏰 Encounter Start")
	return sb.String()
}

func printRound(r event.Round) string {
	sb := strings.Builder{}
	sb.WriteString("🔄 Round ")
	sb.WriteString(fmt.Sprint(r.Round))
	return sb.String()
}

func printTurn(t event.Turn) string {
	sb := strings.Builder{}
	sb.WriteString("🔃 Turn ")
	sb.WriteString(fmt.Sprint(t.Turn))
	sb.WriteString(": ")
	sb.WriteString(fmt.Sprint(t.Creature.Name))
	return sb.String()
}

func printDeath(d event.Death) string {
	sb := strings.Builder{}
	sb.WriteString("☠️ ")
	sb.WriteString(fmt.Sprint(d.Creature.Name))
	sb.WriteString(" is about to die")
	return sb.String()
}

func printUseAction(u event.UseAction) string {
	sb := strings.Builder{}
	sb.WriteString("🗡️ ")
	sb.WriteString(fmt.Sprint(u.Source.Name))
	sb.WriteString(" uses ")
	sb.WriteString(fmt.Sprint(u.Action.Name))
	sb.WriteString(" on ")
	sb.WriteString(fmt.Sprint(u.Target.Name))
	return sb.String()
}

func printTakeDamage(d event.TakeDamage) string {
	sb := strings.Builder{}
	sb.WriteString("🩸 ")
	sb.WriteString(fmt.Sprint(d.Creature.Name))
	sb.WriteString(" takes ")
	sb.WriteString(fmt.Sprint(d.Damage))
	sb.WriteString(" damage")
	return sb.String()
}
