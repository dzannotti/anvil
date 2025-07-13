package expression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"anvil/internal/expression"
	"anvil/internal/tag"
)

type mockRoller struct {
	values []int
	index  int
}

func newMockRoller(values ...int) *mockRoller {
	return &mockRoller{values: values}
}

func (m *mockRoller) Roll(_ int) int {
	if m.index >= len(m.values) {
		return 1
	}

	value := m.values[m.index]
	m.index++
	return value
}

func TestFromConstant(t *testing.T) {
	t.Run("creates expression with constant value", func(t *testing.T) {
		expr := expression.FromConstant(5, "test")
		require.NotNil(t, expr)
		assert.Equal(t, 0, expr.Value)
		assert.Len(t, expr.Components, 1)
	})

	t.Run("evaluates to constant value", func(t *testing.T) {
		expr := expression.FromConstant(10, "test")
		expr.Evaluate()
		assert.Equal(t, 10, expr.Value)
	})
}

func TestFromDice(t *testing.T) {
	t.Run("creates expression with dice component", func(t *testing.T) {
		expr := expression.FromDice(2, 6, "test")
		require.NotNil(t, expr)
		assert.Equal(t, 0, expr.Value)
		assert.Len(t, expr.Components, 1)
	})

	t.Run("evaluates dice rolls", func(t *testing.T) {
		expr := expression.FromDice(2, 6, "test")
		expr.Rng = newMockRoller(3, 4)
		expr.Evaluate()
		assert.Equal(t, 7, expr.Value)
	})
}

func TestFromD20(t *testing.T) {
	t.Run("creates expression with d20 component", func(t *testing.T) {
		expr := expression.FromD20("test")
		require.NotNil(t, expr)
		assert.Equal(t, 0, expr.Value)
		assert.Len(t, expr.Components, 1)
	})

	t.Run("evaluates d20 roll", func(t *testing.T) {
		expr := expression.FromD20("test")
		expr.Rng = newMockRoller(15)
		expr.Evaluate()
		assert.Equal(t, 15, expr.Value)
	})
}

func TestFromDamageConstant(t *testing.T) {
	t.Run("creates expression with damage constant", func(t *testing.T) {
		tags := tag.ContainerFromString("damage.fire")
		expr := expression.FromDamageConstant(8, tags, "test")
		require.NotNil(t, expr)
		assert.Equal(t, 0, expr.Value)
		assert.Len(t, expr.Components, 1)
	})

	t.Run("evaluates to damage value", func(t *testing.T) {
		tags := tag.ContainerFromString("damage.fire")
		expr := expression.FromDamageConstant(12, tags, "test")
		expr.Evaluate()
		assert.Equal(t, 12, expr.Value)
	})
}

func TestFromDamageDice(t *testing.T) {
	t.Run("creates expression with damage dice", func(t *testing.T) {
		tags := tag.ContainerFromString("damage.fire")
		expr := expression.FromDamageDice(3, 6, tags, "test")
		require.NotNil(t, expr)
		assert.Equal(t, 0, expr.Value)
		assert.Len(t, expr.Components, 1)
	})

	t.Run("evaluates damage dice", func(t *testing.T) {
		tags := tag.ContainerFromString("damage.fire")
		expr := expression.FromDamageDice(2, 8, tags, "test")
		expr.Rng = newMockRoller(5, 7)
		expr.Evaluate()
		assert.Equal(t, 12, expr.Value)
	})
}

func TestExpression_AddConstant(t *testing.T) {
	t.Run("adds constant to existing expression", func(t *testing.T) {
		expr := expression.FromConstant(5, "base")
		expr.AddConstant(3, "bonus")
		assert.Len(t, expr.Components, 2)

		expr.Evaluate()
		assert.Equal(t, 8, expr.Value)
	})
}

func TestExpression_AddDice(t *testing.T) {
	t.Run("adds dice to existing expression", func(t *testing.T) {
		expr := expression.FromConstant(5, "base")
		expr.AddDice(1, 6, "bonus")
		assert.Len(t, expr.Components, 2)

		expr.Rng = newMockRoller(4)
		expr.Evaluate()
		assert.Equal(t, 9, expr.Value)
	})
}

func TestExpression_AddD20(t *testing.T) {
	t.Run("adds d20 to existing expression", func(t *testing.T) {
		expr := expression.FromConstant(5, "base")
		expr.AddD20("ability")
		assert.Len(t, expr.Components, 2)

		expr.Rng = newMockRoller(12)
		expr.Evaluate()
		assert.Equal(t, 17, expr.Value)
	})
}

func TestExpression_AddDamageConstant(t *testing.T) {
	t.Run("adds damage constant to existing expression", func(t *testing.T) {
		expr := expression.FromConstant(5, "base")
		tags := tag.ContainerFromString("damage.fire")
		expr.AddDamageConstant(3, tags, "fire damage")
		assert.Len(t, expr.Components, 2)

		expr.Evaluate()
		assert.Equal(t, 8, expr.Value)
	})
}

func TestExpression_AddDamageDice(t *testing.T) {
	t.Run("adds damage dice to existing expression", func(t *testing.T) {
		expr := expression.FromConstant(5, "base")
		tags := tag.ContainerFromString("damage.fire")
		expr.AddDamageDice(2, 4, tags, "fire damage")
		assert.Len(t, expr.Components, 2)

		expr.Rng = newMockRoller(3, 4)
		expr.Evaluate()
		assert.Equal(t, 12, expr.Value)
	})
}

func TestExpression_Evaluate(t *testing.T) {
	t.Run("evaluates complex expression", func(t *testing.T) {
		expr := expression.FromConstant(10, "base")
		expr.AddDice(2, 6, "bonus dice")
		expr.AddConstant(5, "flat bonus")
		expr.Rng = newMockRoller(3, 4)

		result := expr.Evaluate()
		assert.Equal(t, expr, result)
		assert.Equal(t, 22, expr.Value)
	})

	t.Run("re-evaluates expression correctly", func(t *testing.T) {
		expr := expression.FromDice(1, 6, "test")
		expr.Rng = newMockRoller(3, 5)

		expr.Evaluate()
		assert.Equal(t, 3, expr.Value)

		expr.Rng = newMockRoller(5)
		expr.Evaluate()
		assert.Equal(t, 5, expr.Value)
	})
}

func TestExpression_GiveAdvantage(t *testing.T) {
	t.Run("gives advantage to d20 component", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(8, 15)

		expr.GiveAdvantage("source")
		expr.Evaluate()
		assert.Equal(t, 15, expr.Value)
	})

	t.Run("panics when no components", func(t *testing.T) {
		expr := &expression.Expression{}
		assert.Panics(t, func() {
			expr.GiveAdvantage("source")
		})
	})

	t.Run("panics when first component is not d20", func(t *testing.T) {
		expr := expression.FromConstant(5, "test")
		assert.Panics(t, func() {
			expr.GiveAdvantage("source")
		})
	})
}

func TestExpression_GiveDisadvantage(t *testing.T) {
	t.Run("gives disadvantage to d20 component", func(t *testing.T) {
		expr := expression.FromD20("attack")
		expr.Rng = newMockRoller(8, 15)

		expr.GiveDisadvantage("source")
		expr.Evaluate()
		assert.Equal(t, 8, expr.Value)
	})

	t.Run("panics when no components", func(t *testing.T) {
		expr := &expression.Expression{}
		assert.Panics(t, func() {
			expr.GiveDisadvantage("source")
		})
	})

	t.Run("panics when first component is not d20", func(t *testing.T) {
		expr := expression.FromConstant(5, "test")
		assert.Panics(t, func() {
			expr.GiveDisadvantage("source")
		})
	})
}

func TestExpression_ReplaceWith(t *testing.T) {
	t.Run("replaces all components with constant", func(t *testing.T) {
		expr := expression.FromDice(2, 6, "original")
		expr.AddConstant(5, "bonus")
		assert.Len(t, expr.Components, 2)

		tags := tag.ContainerFromString("replaced")
		expr.ReplaceWith(20, "critical", tags)
		assert.Len(t, expr.Components, 1)

		expr.Evaluate()
		assert.Equal(t, 20, expr.Value)
	})
}

func TestExpression_DoubleDice(t *testing.T) {
	t.Run("doubles only dice components", func(t *testing.T) {
		expr := expression.FromDice(2, 6, "weapon")
		expr.AddConstant(5, "modifier")
		assert.Len(t, expr.Components, 2)

		expr.DoubleDice("critical hit")
		assert.Len(t, expr.Components, 3)

		expr.Rng = newMockRoller(3, 4, 5, 6)
		expr.Evaluate()
		assert.Equal(t, 23, expr.Value)
	})

	t.Run("ignores non-dice components", func(t *testing.T) {
		expr := expression.FromConstant(10, "base")
		expr.DoubleDice("critical")
		assert.Len(t, expr.Components, 1)

		expr.Evaluate()
		assert.Equal(t, 10, expr.Value)
	})

	t.Run("handles multiple dice components", func(t *testing.T) {
		expr := expression.FromDice(1, 8, "weapon")
		expr.AddDice(1, 6, "sneak attack")
		assert.Len(t, expr.Components, 2)

		expr.DoubleDice("critical")
		assert.Len(t, expr.Components, 4)

		expr.Rng = newMockRoller(5, 3, 6, 4)
		expr.Evaluate()
		assert.Equal(t, 18, expr.Value)
	})
}

func TestExpression_MaxDice(t *testing.T) {
	t.Run("maximizes only dice components", func(t *testing.T) {
		expr := expression.FromDice(2, 6, "weapon")
		expr.AddConstant(3, "modifier")
		assert.Len(t, expr.Components, 2)

		expr.MaxDice("brutal critical")
		assert.Len(t, expr.Components, 3)

		expr.Rng = newMockRoller(1, 2)
		expr.Evaluate()
		assert.Equal(t, 18, expr.Value)
	})

	t.Run("ignores non-dice components", func(t *testing.T) {
		expr := expression.FromConstant(10, "base")
		expr.MaxDice("brutal")
		assert.Len(t, expr.Components, 1)

		expr.Evaluate()
		assert.Equal(t, 10, expr.Value)
	})

	t.Run("handles multiple dice components", func(t *testing.T) {
		expr := expression.FromDice(2, 8, "weapon")
		expr.AddDice(3, 4, "poison")
		assert.Len(t, expr.Components, 2)

		expr.MaxDice("brutal critical")
		assert.Len(t, expr.Components, 4)

		expr.Rng = newMockRoller(1, 1, 1, 1, 1)
		expr.Evaluate()
		assert.Equal(t, 33, expr.Value)
	})
}

func TestExpression_EvaluateDamage(t *testing.T) {
	t.Run("returns empty expression when no components", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		result := expr.EvaluateDamage()

		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Value)
		assert.Len(t, result.Components, 0)
	})

	t.Run("groups damage by tag types", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(5, tag.NewContainer(expression.DamageSlashing), "sword")
		expr.AddDamageConstant(3, tag.NewContainer(expression.DamageFire), "fire enchant")
		expr.AddDamageConstant(2, tag.NewContainer(expression.DamageFire), "more fire")

		result := expr.EvaluateDamage()
		assert.Equal(t, 10, result.Value)
		assert.Len(t, result.Components, 2)

		slashingComp := findComponentByTag(result, expression.DamageSlashing)
		fireComp := findComponentByTag(result, expression.DamageFire)

		assert.NotNil(t, slashingComp)
		assert.NotNil(t, fireComp)
		assert.Equal(t, 5, slashingComp.Value())
		assert.Equal(t, 5, fireComp.Value())
	})

	t.Run("groups primary tags under first component tags", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(8, tag.NewContainer(expression.DamageSlashing), "base weapon")
		expr.AddConstant(3, "strength modifier")
		expr.AddConstant(2, "magic bonus")

		result := expr.EvaluateDamage()
		assert.Equal(t, 13, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		assert.Equal(t, 13, comp.Value())
		tags := comp.Tags()
		assert.True(t, tags.HasTag(expression.DamageSlashing))
	})

	t.Run("handles mixed primary and typed damage", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(6, tag.NewContainer(expression.DamagePiercing), "weapon")
		expr.AddConstant(4, "strength")
		expr.AddDamageConstant(3, tag.NewContainer(expression.DamagePoison), "poison")
		expr.AddConstant(1, "magic")

		result := expr.EvaluateDamage()
		assert.Equal(t, 14, result.Value)
		assert.Len(t, result.Components, 2)

		piercingComp := findComponentByTag(result, expression.DamagePiercing)
		poisonComp := findComponentByTag(result, expression.DamagePoison)

		assert.NotNil(t, piercingComp)
		assert.NotNil(t, poisonComp)
		assert.Equal(t, 11, piercingComp.Value())
		assert.Equal(t, 3, poisonComp.Value())
	})

	t.Run("handles dice components in grouping", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller(4, 2, 6)}
		expr.AddDamageDice(1, 8, tag.NewContainer(expression.DamageSlashing), "weapon")
		expr.AddDamageDice(1, 6, tag.NewContainer(expression.DamageSlashing), "sneak attack")
		expr.AddDamageDice(1, 4, tag.NewContainer(expression.DamageFire), "fire")

		result := expr.EvaluateDamage()
		assert.Equal(t, 12, result.Value)
		assert.Len(t, result.Components, 2)

		slashingComp := findComponentByTag(result, expression.DamageSlashing)
		fireComp := findComponentByTag(result, expression.DamageFire)

		assert.NotNil(t, slashingComp)
		assert.NotNil(t, fireComp)
		assert.Equal(t, 6, slashingComp.Value())
		assert.Equal(t, 6, fireComp.Value())
	})

	t.Run("preserves first component tags for primary grouping", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(10, tag.NewContainer(expression.DamageForce), "magic missile")
		expr.AddConstant(5, "spell modifier")

		result := expr.EvaluateDamage()
		assert.Equal(t, 15, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		tags := comp.Tags()
		assert.True(t, tags.HasTag(expression.DamageForce))
		assert.Equal(t, "magic missile", comp.Source())
	})

	t.Run("handles empty tag components", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(5, tag.NewContainer(expression.DamageSlashing), "weapon")
		emptyExpr := expression.FromConstant(3, "empty tag component")
		expr.Components = append(expr.Components, emptyExpr.Components[0])

		result := expr.EvaluateDamage()
		assert.Equal(t, 8, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		tags := comp.Tags()
		assert.True(t, tags.HasTag(expression.DamageSlashing))
		assert.Equal(t, 8, comp.Value())
	})

	t.Run("builds appropriate group sources", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(3, tag.NewContainer(expression.DamageFire), "fire spell")
		expr.AddDamageConstant(2, tag.NewContainer(expression.DamageFire), "fire enchant")
		expr.AddDamageConstant(1, tag.NewContainer(expression.DamageFire), "fire aura")

		result := expr.EvaluateDamage()
		assert.Equal(t, 6, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		assert.Equal(t, "grouped damage (3 sources)", comp.Source())
	})

	t.Run("handles edge case with empty groups", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		result := expr.EvaluateDamage()
		assert.Equal(t, 0, result.Value)
		assert.Len(t, result.Components, 0)
	})

	t.Run("hits groupComponentsByTags early return with empty expression", func(t *testing.T) {
		// Create completely empty expression (no components at all)
		expr := &expression.Expression{Rng: newMockRoller()}
		// This should hit line 178: early return when len(e.Components) == 0
		result := expr.EvaluateDamage()
		assert.Equal(t, 0, result.Value)
		assert.Len(t, result.Components, 0)
	})

	t.Run("handles expression with only empty tag components", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		emptyExpr1 := expression.FromConstant(5, "first")
		emptyExpr2 := expression.FromConstant(3, "second")
		expr.Components = append(expr.Components, emptyExpr1.Components[0], emptyExpr2.Components[0])

		result := expr.EvaluateDamage()
		assert.Equal(t, 8, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		assert.Equal(t, 8, comp.Value())
		assert.Equal(t, "grouped damage (2 sources)", comp.Source())
	})

	t.Run("preserves single component source when same", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		expr.AddDamageConstant(5, tag.NewContainer(expression.DamageFire), "fire spell")

		result := expr.EvaluateDamage()
		assert.Equal(t, 5, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		assert.Equal(t, "fire spell", comp.Source())
	})

	t.Run("preserves same source across multiple components in group", func(t *testing.T) {
		expr := &expression.Expression{Rng: newMockRoller()}
		// Add multiple components with same source and same tag
		expr.AddDamageConstant(3, tag.NewContainer(expression.DamageFire), "same source")
		expr.AddDamageConstant(2, tag.NewContainer(expression.DamageFire), "same source")
		expr.AddDamageConstant(1, tag.NewContainer(expression.DamageFire), "same source")

		result := expr.EvaluateDamage()
		assert.Equal(t, 6, result.Value)
		assert.Len(t, result.Components, 1)

		comp := result.Components[0]
		// This should hit line 246 - the fallback return for same sources
		assert.Equal(t, "same source", comp.Source())
	})
}

func TestDamageConstants(t *testing.T) {
	t.Run("damage constants are properly defined", func(t *testing.T) {
		constants := map[string]tag.Tag{
			"damage.acid":        expression.DamageAcid,
			"damage.bludgeoning": expression.DamageBludgeoning,
			"damage.cold":        expression.DamageCold,
			"damage.fire":        expression.DamageFire,
			"damage.force":       expression.DamageForce,
			"damage.lightning":   expression.DamageLightning,
			"damage.necrotic":    expression.DamageNecrotic,
			"damage.piercing":    expression.DamagePiercing,
			"damage.poison":      expression.DamagePoison,
			"damage.psychic":     expression.DamagePsychic,
			"damage.radiant":     expression.DamageRadiant,
			"damage.slashing":    expression.DamageSlashing,
			"damage.thunder":     expression.DamageThunder,
		}

		for expected, actual := range constants {
			assert.Equal(t, expected, actual.AsString())
			assert.True(t, actual.IsValid())
		}
	})
}

func findComponentByTag(expr *expression.Expression, targetTag tag.Tag) expression.Component {
	for _, comp := range expr.Components {
		tags := comp.Tags()
		if tags.HasTag(targetTag) {
			return comp
		}
	}
	return nil
}
