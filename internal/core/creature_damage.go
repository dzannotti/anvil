package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(c.hitPoints-damage, 0)
	c.log.Add(NewTakeDamageEvent(c, damage))
}

func (c *Creature) StartTurn() {

}

func (c *Creature) AttackRoll(target definition.Creature, tc tag.Container) definition.CheckResult {
	expression := expression.FromD20("Base")
	c.log.Start(NewAttackRollEvent(c, target))
	defer c.log.End()
	before := BeforeAttackRollState{Source: c, Target: target, Expression: &expression, Tags: tc}
	c.effects.Evaluate("BeforeAttackRoll", before)
	expression.Evaluate()
	after := AfterAttackRollState{Source: c, Target: target, Result: &expression, Tags: tc}
	c.effects.Evaluate("AfterAttackRoll", after)
	c.log.Add(NewExpressionResultEvent(expression))
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	c.log.Add(NewAttributeCalculationEvent(tags.ArmorClass, targetAC))
	ok := value >= targetAC.Value
	c.log.Add(NewCheckResultEvent(value, targetAC.Value, crit, ok))
	return definition.NewCheckResult(value, expression, crit, ok)
}
