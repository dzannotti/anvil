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
	expression := expression.FromD20("Base")
	c.Log.Start(AttackRollEventType, AttackRollEvent{Source: *c, Target: *target})
	defer c.Log.End()
	before := BeforeAttackRollState{Source: c, Target: target, Expression: &expression, Tags: tc}
	c.Effects.Evaluate(BeforeAttackRollStateType, before)
	expression.Evaluate()
	after := AfterAttackRollState{Source: c, Target: target, Result: &expression, Tags: tc}
	c.Effects.Evaluate(AfterAttackRollStateType, after)
	c.Log.Add(ExpressionResultEventType, ExpressionResultEvent{Expression: expression})
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	c.Log.Add(AttributeCalculationEventType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	c.Log.Add(CheckResultEventType, CheckResultEvent{Value: value, Against: targetAC.Value, Critical: crit, Success: ok})
	return CheckResult{Value: value, Against: expression.Value, Critical: crit, Success: ok}
}
