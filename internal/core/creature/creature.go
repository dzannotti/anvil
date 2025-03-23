package creature

import (
	"anvil/internal/core/team"
	"anvil/internal/log"
	"math/rand"
)

type Creature struct {
	name         string
	team         team.Team
	hitPoints    int
	maxHitPoints int
	actionPoints int
	actions      []Action
	log          *log.EventLog
}

func (c Creature) Name() string {
	return c.name
}

func (c Creature) HitPoints() int {
	return c.hitPoints
}

func (c Creature) MaxHitPoints() int {
	return c.maxHitPoints
}

func (c Creature) Team() team.Team {
	return c.team
}

func (c Creature) ActionPoints() int {
	return c.actionPoints
}

func RollDice(sides int) int {
	return rand.Intn(sides) + 1
}

func New(log *log.EventLog, name string, t team.Team, hp int, actions []Action) *Creature {
	return &Creature{log: log, name: name, team: t, hitPoints: hp, actionPoints: 0, actions: actions, maxHitPoints: hp}
}

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(0, c.hitPoints-damage)
	c.log.Start(NewTakeDamageEvent(c, damage))
	c.log.End()
}

func (c *Creature) Attack(target *Creature) {
	c.log.Start(NewUseActionEvent("Attack", c, target))
	damage := RollDice(20)
	target.TakeDamage(damage)
	if target.IsDead() {
		target.Die()
	}
	c.log.End()
}

func (c *Creature) IsDead() bool {
	return c.hitPoints == 0
}

func (c *Creature) StartTurn() {
	c.actionPoints = 1
}

func (c *Creature) Consume(cost int) {
	c.actionPoints = max(0, c.actionPoints-cost)
}

func (c *Creature) Actions() []Action {
	return c.actions
}

func (c *Creature) Die() {
	c.log.Start(NewDeathEvent(c))
	c.log.End()
}
