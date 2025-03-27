package prettyprint

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"anvil/internal/core/event"
	"anvil/internal/core/event/snapshot"
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
	}
	return "unknown event" + reflect.TypeOf(ev.Data).Name()
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

	return sb.String()
}

func printRound(r event.Round) string {
	sb := strings.Builder{}
	sb.WriteString("🔄 Round ")
	sb.WriteString(fmt.Sprint(r.Round + 1))
	return sb.String()
}

func printTurn(t event.Turn) string {
	sb := strings.Builder{}
	sb.WriteString("🔃 Turn ")
	sb.WriteString(fmt.Sprint(t.Turn + 1))
	sb.WriteString(": ")
	sb.WriteString(fmt.Sprint(t.Creature.Name))
	return sb.String()
}

func printDeath(d event.Died) string {
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
	sb.WriteString(fmt.Sprint(d.Target.Name))
	sb.WriteString(" takes ")
	sb.WriteString(fmt.Sprint(d.Damage))
	sb.WriteString(" damage")
	return sb.String()
}
