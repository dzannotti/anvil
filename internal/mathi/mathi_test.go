package mathi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	t.Run("positive number returns itself", func(t *testing.T) {
		assert.Equal(t, 5, Abs(5))
		assert.Equal(t, 10, Abs(10))
		assert.Equal(t, 1, Abs(1))
	})

	t.Run("negative number returns positive", func(t *testing.T) {
		assert.Equal(t, 5, Abs(-5))
		assert.Equal(t, 10, Abs(-10))
		assert.Equal(t, 1, Abs(-1))
	})

	t.Run("zero returns zero", func(t *testing.T) {
		assert.Equal(t, 0, Abs(0))
	})
}

func TestMax(t *testing.T) {
	t.Run("returns larger of two positive numbers", func(t *testing.T) {
		assert.Equal(t, 10, Max(5, 10))
		assert.Equal(t, 10, Max(10, 5))
		assert.Equal(t, 7, Max(7, 3))
	})

	t.Run("returns larger of two negative numbers", func(t *testing.T) {
		assert.Equal(t, -3, Max(-5, -3))
		assert.Equal(t, -1, Max(-1, -10))
	})

	t.Run("handles mixed positive and negative", func(t *testing.T) {
		assert.Equal(t, 5, Max(-3, 5))
		assert.Equal(t, 0, Max(-5, 0))
		assert.Equal(t, 1, Max(1, -1))
	})

	t.Run("equal values return either", func(t *testing.T) {
		assert.Equal(t, 5, Max(5, 5))
		assert.Equal(t, 0, Max(0, 0))
		assert.Equal(t, -3, Max(-3, -3))
	})
}

func TestMin(t *testing.T) {
	t.Run("returns smaller of two positive numbers", func(t *testing.T) {
		assert.Equal(t, 5, Min(5, 10))
		assert.Equal(t, 5, Min(10, 5))
		assert.Equal(t, 3, Min(7, 3))
	})

	t.Run("returns smaller of two negative numbers", func(t *testing.T) {
		assert.Equal(t, -5, Min(-5, -3))
		assert.Equal(t, -10, Min(-1, -10))
	})

	t.Run("handles mixed positive and negative", func(t *testing.T) {
		assert.Equal(t, -3, Min(-3, 5))
		assert.Equal(t, -5, Min(-5, 0))
		assert.Equal(t, -1, Min(1, -1))
	})

	t.Run("equal values return either", func(t *testing.T) {
		assert.Equal(t, 5, Min(5, 5))
		assert.Equal(t, 0, Min(0, 0))
		assert.Equal(t, -3, Min(-3, -3))
	})
}

func TestClamp(t *testing.T) {
	t.Run("value within range returns unchanged", func(t *testing.T) {
		assert.Equal(t, 5, Clamp(5, 0, 10))
		assert.Equal(t, 7, Clamp(7, 5, 10))
		assert.Equal(t, 0, Clamp(0, -5, 5))
	})

	t.Run("value above range returns high bound", func(t *testing.T) {
		assert.Equal(t, 10, Clamp(15, 0, 10))
		assert.Equal(t, 5, Clamp(100, -10, 5))
		assert.Equal(t, 0, Clamp(1, -5, 0))
	})

	t.Run("value below range returns low bound", func(t *testing.T) {
		assert.Equal(t, 0, Clamp(-5, 0, 10))
		assert.Equal(t, -10, Clamp(-100, -10, 5))
		assert.Equal(t, -5, Clamp(-10, -5, 0))
	})

	t.Run("value equals bounds", func(t *testing.T) {
		assert.Equal(t, 0, Clamp(0, 0, 10))
		assert.Equal(t, 10, Clamp(10, 0, 10))
		assert.Equal(t, 5, Clamp(5, 5, 5))
	})

	t.Run("handles negative ranges", func(t *testing.T) {
		assert.Equal(t, -3, Clamp(-1, -5, -3))
		assert.Equal(t, -5, Clamp(-10, -5, -3))
		assert.Equal(t, -4, Clamp(-4, -5, -3))
	})
}

func TestSign(t *testing.T) {
	t.Run("positive numbers return 1", func(t *testing.T) {
		assert.Equal(t, 1, Sign(1))
		assert.Equal(t, 1, Sign(5))
		assert.Equal(t, 1, Sign(100))
	})

	t.Run("negative numbers return -1", func(t *testing.T) {
		assert.Equal(t, -1, Sign(-1))
		assert.Equal(t, -1, Sign(-5))
		assert.Equal(t, -1, Sign(-100))
	})

	t.Run("zero returns 1", func(t *testing.T) {
		assert.Equal(t, 1, Sign(0))
	})
}
