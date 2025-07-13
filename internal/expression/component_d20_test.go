package expression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/expression"
)

func TestD20Component_Evaluate(t *testing.T) {
	t.Run("rolls single d20 when no advantage or disadvantage", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(12)
		expr.Evaluate()

		assert.Equal(t, 12, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, expression.ComponentKindD20, comp.Kind())
		assert.Equal(t, 12, comp.Value())
		assert.Equal(t, "attack", comp.Source())
	})

	t.Run("rolls with advantage takes higher value", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(8, 15)
		expr.GiveAdvantage("high ground")
		expr.Evaluate()

		assert.Equal(t, 15, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, 15, comp.Value())
	})

	t.Run("rolls with disadvantage takes lower value", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(8, 15)
		expr.GiveDisadvantage("prone")
		expr.Evaluate()

		assert.Equal(t, 8, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, 8, comp.Value())
	})

	t.Run("advantage and disadvantage cancel out", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(12)
		expr.GiveAdvantage("blessing")
		expr.GiveDisadvantage("cursed")
		expr.Evaluate()

		assert.Equal(t, 12, expr.Value)
		comp := expr.Components[0]
		assert.Equal(t, 12, comp.Value())
	})

	t.Run("multiple advantages still just advantage", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(5, 18)
		expr.GiveAdvantage("blessing")
		expr.GiveAdvantage("luck")
		expr.Evaluate()

		assert.Equal(t, 18, expr.Value)
	})

	t.Run("multiple disadvantages still just disadvantage", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(5, 18)
		expr.GiveDisadvantage("exhaustion")
		expr.GiveDisadvantage("blinded")
		expr.Evaluate()

		assert.Equal(t, 5, expr.Value)
	})

	t.Run("multiple advantages + one disadvantage cancel out (D&D 5e rule)", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(12)
		expr.GiveAdvantage("blessing")
		expr.GiveAdvantage("high ground")
		expr.GiveAdvantage("flanking")
		expr.GiveDisadvantage("prone target")
		expr.Evaluate()

		assert.Equal(t, 12, expr.Value)
	})

	t.Run("multiple disadvantages + one advantage cancel out (D&D 5e rule)", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(15)
		expr.GiveDisadvantage("exhaustion")
		expr.GiveDisadvantage("blinded")
		expr.GiveDisadvantage("restrained")
		expr.GiveAdvantage("blessing")
		expr.Evaluate()

		assert.Equal(t, 15, expr.Value)
	})

	t.Run("2 advantages + 2 disadvantages cancel out", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(8)
		expr.GiveAdvantage("blessing")
		expr.GiveAdvantage("flanking")
		expr.GiveDisadvantage("exhaustion")
		expr.GiveDisadvantage("blinded")
		expr.Evaluate()

		assert.Equal(t, 8, expr.Value)
	})
}

func TestD20Component_CriticalChecks(t *testing.T) {
	t.Run("detects critical success", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(20)
		expr.Evaluate()

		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)
		assert.True(t, d20Comp.IsCriticalSuccess())
		assert.False(t, d20Comp.IsCriticalFailure())
		assert.True(t, d20Comp.IsCritical())
	})

	t.Run("detects critical failure", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(1)
		expr.Evaluate()

		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)
		assert.False(t, d20Comp.IsCriticalSuccess())
		assert.True(t, d20Comp.IsCriticalFailure())
		assert.True(t, d20Comp.IsCritical())
	})

	t.Run("detects non-critical roll", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(10)
		expr.Evaluate()

		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)
		assert.False(t, d20Comp.IsCriticalSuccess())
		assert.False(t, d20Comp.IsCriticalFailure())
		assert.False(t, d20Comp.IsCritical())
	})

	t.Run("critical with advantage uses final result", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(10, 20)
		expr.GiveAdvantage("luck")
		expr.Evaluate()

		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)
		assert.True(t, d20Comp.IsCriticalSuccess())
		assert.Equal(t, 20, d20Comp.Value())
	})

	t.Run("critical with disadvantage uses final result", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(20, 1)
		expr.GiveDisadvantage("cursed")
		expr.Evaluate()

		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)
		assert.True(t, d20Comp.IsCriticalFailure())
		assert.Equal(t, 1, d20Comp.Value())
	})
}

func TestD20Component_GiveAdvantage(t *testing.T) {
	t.Run("can be given advantage directly", func(t *testing.T) {
		expr := expression.FromD20("attack")
		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)

		result := d20Comp.GiveAdvantage("blessing")
		assert.Equal(t, d20Comp, result)

		expr.Rng = newMockRoller(5, 18)
		expr.Evaluate()
		assert.Equal(t, 18, expr.Value)
	})
}

func TestD20Component_GiveDisadvantage(t *testing.T) {
	t.Run("can be given disadvantage directly", func(t *testing.T) {
		expr := expression.FromD20("attack")
		comp := expr.Components[0]
		d20Comp, ok := comp.(*expression.D20Component)
		assert.True(t, ok)

		result := d20Comp.GiveDisadvantage("exhaustion")
		assert.Equal(t, d20Comp, result)

		expr.Rng = newMockRoller(5, 18)
		expr.Evaluate()
		assert.Equal(t, 5, expr.Value)
	})
}
