package sendclientmessagejob_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sendclientmessagejob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func TestMarshalPayload_Smoke(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := sendclientmessagejob.MarshalPayload(types.NewMessageID())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := sendclientmessagejob.MarshalPayload(types.MessageIDNil)
		require.Error(t, err)
		assert.Empty(t, p)
	})
}

func TestUnmarshalPayload(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := sendclientmessagejob.UnmarshalPayload(types.NewMessageID().String())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := sendclientmessagejob.UnmarshalPayload("invalid payload")
		require.Error(t, err)
		assert.Empty(t, p)
	})

	t.Run("empty payload", func(t *testing.T) {
		p, err := sendclientmessagejob.UnmarshalPayload("")
		require.Error(t, err)
		assert.Empty(t, p)
	})

	t.Run("empty id", func(t *testing.T) {
		p, err := sendclientmessagejob.UnmarshalPayload(types.MessageIDNil.String())
		require.Error(t, err)
		assert.Empty(t, p)
	})
}
