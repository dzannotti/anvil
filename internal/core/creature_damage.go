package core

import (
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func (c *Creature) TakeDamage(damage int) {
	c.HitPoints = max(c.HitPoints-damage, 0)
	c.Log.Add(TakeDamageEventType, TakeDamageEvent{Target: *c, Damage: damage})
}

func (c *Creature) StartTurn() {

}

func (c *Creature) AttackRoll(target *Creature, tc tag.Container) CheckResult {
	expr := expression.FromD20("Base")
	c.Log.Start(AttackRollEventType, AttackRollEvent{Source: *c, Target: *target})
	defer c.Log.End()
	before := BeforeAttackRollState{Source: c, Target: target, Expression: &expr, Tags: tc}
	c.Effects.Evaluate(BeforeAttackRollStateType, before)
	expr.Evaluate()
	after := AfterAttackRollState{Source: c, Target: target, Result: &expr, Tags: tc}
	c.Effects.Evaluate(AfterAttackRollStateType, after)
	c.Log.Add(ExpressionResultEventType, ExpressionResultEvent{Expression: expr})
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	c.Log.Add(AttributeCalculationEventType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	c.Log.Add(CheckResultEventType, CheckResultEvent{Value: value, Against: targetAC.Value, Critical: crit, Success: ok})
	return CheckResult{Value: value, Against: expr.Value, Critical: crit, Success: ok}
}
