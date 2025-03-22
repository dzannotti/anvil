package creature

import (
	"anvil/internal/team"
	"fmt"
	"math/rand"
)

type Creature struct {
	name         string
	factionID    team.Team
	hitPoints    int
	actionPoints int
	actions      []Action
}

func (c Creature) Name() string {
	return c.name
}

func (c Creature) FactionID() team.Team {
	return c.factionID
}

func (c Creature) HitPoints() int {
	return c.hitPoints
}

func (c Creature) ActionPoints() int {
	return c.actionPoints
}

func RollDice(sides int) int {
	return rand.Intn(sides) + 1
}

func New(name string, factionID team.Team, hitPoints int, actions []Action) *Creature {
	return &Creature{name: name, factionID: factionID, hitPoints: hitPoints, actionPoints: 0, actions: actions}
}

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(0, c.hitPoints-damage)
	fmt.Println(c.name, "took", damage, "damage", c.hitPoints, "remaining")
}

func (c *Creature) Attack(target *Creature) {
	fmt.Println(c.name, "attacks", target.name)
	damage := RollDice(20)
	target.TakeDamage(damage)
	if target.IsDead() {
		fmt.Println(target.name, "is dead")
	}
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
