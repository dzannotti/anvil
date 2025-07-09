package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestManager(t *testing.T) {
	t.Run("Basic Request Handling", func(t *testing.T) {
		t.Run("should handle a basic request", func(t *testing.T) {
			rm := NewRequestManager()
			actor := &Actor{Name: "TestActor"}
			options := []RequestOption{
				{Label: "Option 1", Value: "value1"},
				{Label: "Option 2", Value: "value2"},
			}

			// Start request in goroutine since Ask blocks
			go func() {
				request := rm.GetPendingRequest()
				require.NotNil(t, request)
				assert.Equal(t, "TestActor", request.Target.Name)
				assert.Equal(t, "Test question", request.Text)
				assert.Equal(t, 2, len(request.Options))
				request.Answer(options[0])
			}()

			result, err := rm.Ask(actor, "Test question", options)
			require.NoError(t, err)
			assert.Equal(t, "Option 1", result.Label)
			assert.Equal(t, "value1", result.Value)
		})

		t.Run("should return error when already has pending request", func(t *testing.T) {
			rm := NewRequestManager()
			actor := &Actor{Name: "TestActor"}
			options := []RequestOption{{Label: "Option 1", Value: "value1"}}

			// Start first request
			go func() {
				request := rm.GetPendingRequest()
				require.NotNil(t, request)
				request.Answer(options[0])
			}()

			// Start second request before first is answered
			go func() {
				_, err := rm.Ask(actor, "Second question", options)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "already a pending request")
			}()

			// Answer first request
			result, err := rm.Ask(actor, "First question", options)
			require.NoError(t, err)
			assert.Equal(t, "Option 1", result.Label)
		})
	})

	t.Run("Pending Request Checks", func(t *testing.T) {
		t.Run("should return false when no pending request", func(t *testing.T) {
			rm := NewRequestManager()
			assert.False(t, rm.HasPendingRequest())
		})

		t.Run("should return nil when no pending request", func(t *testing.T) {
			rm := NewRequestManager()
			assert.Nil(t, rm.GetPendingRequest())
		})
	})

	t.Run("Answer Default", func(t *testing.T) {
		t.Run("should answer with default option", func(t *testing.T) {
			rm := NewRequestManager()
			actor := &Actor{Name: "TestActor"}
			options := []RequestOption{
				{Label: "Option 1", Value: "value1"},
				{Label: "Option 2", Value: "value2", Default: true},
			}

			go func() {
				for !rm.HasPendingRequest() {
					time.Sleep(time.Millisecond)
				}
				_ = rm.AnswerDefault()
			}()

			result, err := rm.Ask(actor, "Test question", options)
			require.NoError(t, err)
			assert.Equal(t, "Option 2", result.Label)
			assert.Equal(t, "value2", result.Value)
		})

		t.Run("should answer with first option when no default", func(t *testing.T) {
			rm := NewRequestManager()
			actor := &Actor{Name: "TestActor"}
			options := []RequestOption{
				{Label: "Option 1", Value: "value1"},
				{Label: "Option 2", Value: "value2"},
			}

			go func() {
				for !rm.HasPendingRequest() {
					time.Sleep(time.Millisecond)
				}
				_ = rm.AnswerDefault()
			}()

			result, err := rm.Ask(actor, "Test question", options)
			require.NoError(t, err)
			assert.Equal(t, "Option 1", result.Label)
			assert.Equal(t, "value1", result.Value)
		})

		t.Run("should return error when no pending request", func(t *testing.T) {
			rm := NewRequestManager()
			err := rm.AnswerDefault()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "no pending request")
		})
	})
}
