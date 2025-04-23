package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	ui "anvil/cmd/gui/render"
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/demo"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/prettyprint"
)

func printOverhead(ev eventbus.Message, overhead *ui.OverheadManager) {
	var pos grid.Position
	var text string
	color := ui.Text
	switch ev.Kind {
	/*case core.UseActionType:
	data := ev.Data.(core.UseActionEvent)
	pos = data.Source.Position
	text = fmt.Sprintf("*%s*", data.Action.Name())*/
	case core.MoveType:
		data := ev.Data.(core.MoveEvent)
		pos = data.Source.Position
		text = "*move*"
	case core.TakeDamageType:
		data := ev.Data.(core.TakeDamageEvent)
		pos = data.Target.Position
		text = fmt.Sprintf("-%d", data.Damage.Value)
		color = ui.Red
	case core.ConditionChangedType:
		data := ev.Data.(core.ConditionChangedEvent)
		pos = data.Source.Position
		prefix := ""
		color = ui.Yellow
		if !data.Added {
			prefix = "-"
		}
		text = fmt.Sprintf("%s%s", prefix, tags.ToReadableShort(data.Condition))
	case core.EffectType:
		data := ev.Data.(core.EffectEvent)
		pos = data.Source.Position
		color = ui.Yellow
		text = data.Effect.Name
	case core.SavingThrowResultType:
		data := ev.Data.(core.SavingThrowResultEvent)
		pos = data.Actor.Position
		text = "saved"
		color = ui.Green
		if !data.Success {
			text = "failed save"
			color = ui.Yellow
		}
	case core.CheckResultType:
		data := ev.Data.(core.CheckResultEvent)
		if data.Success || !data.Tags.HasTag(tags.Attack) {
			return
		}
		pos = data.Actor.Position
		text = "** miss **"
	default:
		return
	}
	overhead.Add(pos, text, color)
}

func client(_ net.Conn) {
	log := ui.ScrollText{
		Rect:       ui.Rectangle{X: 600, Y: 40, Width: 650, Height: 580},
		LineHeight: 18 + 4,
		Padding:    4,
		BgColor:    ui.Surface0,
		TextColor:  ui.Text,
		FontSize:   18,
	}
	overhead := ui.OverheadManager{}

	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(&log, msg)
		if msg.End {
			return
		}
		printOverhead(msg, &overhead)
	})

	window := ui.Window{}
	window.Open()
	defer window.Close()
	ui.Init()
	defer ui.Close()
	gameState := demo.New(&hub)
	world := gameState.World
	encounter := gameState.Encounter

	camera := ui.Camera{}
	camera.Reset(window.Width, window.Height)
	camera.SetPosition(615, 330)
	am := ui.ActionManager{Encounter: encounter, World: world}

	endTurn := func() {
		encounter.EndTurn()
		am.SetActive(nil)
		if encounter.IsOver() {
			log.AddLine("***** Game over! *****")
			winner, _ := encounter.Winner()
			log.AddLine(fmt.Sprintf("%s won!", string(winner)))
		}
	}

	am.EndTurn = endTurn

	keyBinds := ui.KeyBinds{
		SelectAction: func(i int) {
			actor := encounter.ActiveActor()
			if i > len(actor.Actions) {
				endTurn()
				return
			}
			am.SetActive(actor.Actions[i-1])
		},
	}

	encounter.Start()
	for !window.ShouldClose() {
		dt := ui.GetFrameTime()
		window.StartFrame()
		camera.Begin()
		ui.DrawWorld(world, encounter)
		am.Draw(camera)
		camera.End()
		ui.DrawHeading(encounter)
		ui.DrawActions(encounter.ActiveActor(), am.SetActive, am.Active, endTurn)
		overhead.Draw()
		log.Draw()
		window.EndFrame()
		consumed := ui.ProcessInput()
		if !consumed {
			am.ProcessInput(camera)
		}
		overhead.Update(dt)
		camera.Update()
		ui.Update()
		keyBinds.Update()
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	client(nil)
}
