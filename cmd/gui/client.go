package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"anvil/cmd/gui/render"
	"anvil/cmd/gui/ui"
	"anvil/internal/core"
	"anvil/internal/demo"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/prettyprint"
)

func printOverhead(ev eventbus.Message, overhead *render.OverheadManager) {
	var pos grid.Position
	var text string
	switch ev.Kind {
	case core.UseActionType:
		data := ev.Data.(core.UseActionEvent)
		pos = data.Source.Position
		text = fmt.Sprintf("*%s*", data.Action.Name())
	default:
		return
		/*	case core.TakeDamageType:
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
			case core.SpendResourceType:
				return printSpendResource(ev.Data.(core.SpendResourceEvent))
			case core.ConditionChangedType:
				return printConditionChanged(ev.Data.(core.ConditionChangedEvent))
			case core.MoveType:
				return printMove(ev.Data.(core.MoveEvent))
			case core.MoveStepType:
				return printMoveStep(ev.Data.(core.MoveStepEvent))
			case core.DeathSavingThrowType:
				return printDeathSavingThrow(ev.Data.(core.DeathSavingThrowEvent))
			case core.DeathSavingThrowResultType:
				return printDeathSavingThrowResult(ev.Data.(core.DeathSavingThrowResultEvent))
			case core.DeathSavingThrowAutomaticType:
				return printDeathSavingThrowAutomaticResult(ev.Data.(core.DeathSavingThrowAutomaticEvent))*/
	}
	overhead.Add(pos, text, ui.Red)
}

func client(_ net.Conn) {
	log := ui.ScrollText{
		Rect:       ui.Rectangle{X: 600, Y: 40, Width: 650, Height: 580},
		LineHeight: 18 + 4,
		Padding:    4,
		BgColor:    ui.LightGray,
		TextColor:  ui.Black,
		FontSize:   18,
	}
	overhead := render.OverheadManager{}

	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(&log, msg)
		printOverhead(msg, &overhead)
	})

	window := render.Window{}
	window.Open()
	defer window.Close()
	ui.Init()
	defer ui.Close()
	world, encounter := demo.Create(&hub)
	camera := render.Camera{}
	camera.Reset(window.Width, window.Height)
	camera.SetPosition(-20, -20)
	am := render.ActionManager{Encounter: encounter, World: world}

	endTurn := func() {
		encounter.EndTurn()
		am.SetActive(nil)
	}

	keybinds := ui.KeyBinds{
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
		dt := render.GetFrameTime()
		window.StartFrame()
		camera.Begin()
		render.DrawWorld(world, encounter)
		am.Draw(camera)
		camera.End()
		render.DrawHeading(encounter)
		render.DrawActions(encounter.ActiveActor(), am.SetActive, am.Active, endTurn)
		overhead.Draw(camera)
		log.Draw()
		window.EndFrame()
		consumed := ui.ProcessInput()
		if !consumed {
			am.ProcessInput(camera)
		}
		overhead.Update(dt)
		camera.Update()
		ui.Update()
		keybinds.Update()
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		os.Exit(0)
	}()
	client(nil)
}
