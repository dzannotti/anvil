package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(c.hitPoints-damage, 0)
	c.log.Add(TakeDamageEventType, TakeDamageEvent{Target: c, Damage: damage})
}

func (c *Creature) StartTurn() {

}

func (c *Creature) AttackRoll(target definition.Creature, tc tag.Container) definition.CheckResult {
	expression := expression.FromD20("Base")
	c.log.Start(AttackRollEventType, AttackRollEvent{Source: c, Target: target})
	defer c.log.End()
	before := BeforeAttackRollState{Source: c, Target: target, Expression: &expression, Tags: tc}
	c.effects.Evaluate(BeforeAttackRollStateType, before)
	expression.Evaluate()
	after := AfterAttackRollState{Source: c, Target: target, Result: &expression, Tags: tc}
	c.effects.Evaluate(AfterAttackRollStateType, after)
	c.log.Add(ExpressionResultEventType, ExpressionResultEvent{Expression: expression})
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	c.log.Add(AttributeCalculationEventType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	c.log.Add(CheckResultEventType, CheckResultEvent{Value: value, Against: targetAC.Value, Critical: crit, Success: ok})
	return definition.NewCheckResult(value, expression, crit, ok)
}
