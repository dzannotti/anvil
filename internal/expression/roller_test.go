package expression_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/expression"
)

func TestRngRoller_Roll(t *testing.T) {
	t.Run("rolls within valid range", func(t *testing.T) {
		roller := expression.NewRngRoller()

		for i := 0; i < 100; i++ {
			result := roller.Roll(6)
			assert.GreaterOrEqual(t, result, 1)
			assert.LessOrEqual(t, result, 6)
		}
	})

	t.Run("rolls d20 within valid range", func(t *testing.T) {
		roller := expression.NewRngRoller()

		for i := 0; i < 100; i++ {
			result := roller.Roll(20)
			assert.GreaterOrEqual(t, result, 1)
			assert.LessOrEqual(t, result, 20)
		}
	})

	t.Run("rolls d100 within valid range", func(t *testing.T) {
		roller := expression.NewRngRoller()

		for i := 0; i < 50; i++ {
			result := roller.Roll(100)
			assert.GreaterOrEqual(t, result, 1)
			assert.LessOrEqual(t, result, 100)
		}
	})

	t.Run("different rolls produce different results", func(t *testing.T) {
		roller := expression.NewRngRoller()

		results := make(map[int]bool)
		for i := 0; i < 20; i++ {
			result := roller.Roll(20)
			results[result] = true
		}

		assert.Greater(t, len(results), 1, "Should produce varied results")
	})
}

func TestContext(t *testing.T) {
	t.Run("provides roller to components", func(t *testing.T) {
		mockRoller := newMockRoller(10)
		ctx := &expression.Context{Rng: mockRoller}

		result := ctx.Rng.Roll(20)
		assert.Equal(t, 10, result)
	})
}
