package jobpayload_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/pkg/jobpayload"
)

func TestMarshal_Smoke(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := jobpayload.Marshal(types.NewMessageID())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := jobpayload.Marshal(types.MessageIDNil)
		require.Error(t, err)
		assert.Empty(t, p)
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := jobpayload.Unmarshal(types.NewMessageID().String())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := jobpayload.Unmarshal("invalid payload")
		require.Error(t, err)
		assert.Empty(t, p)
	})

	t.Run("empty payload", func(t *testing.T) {
		p, err := jobpayload.Unmarshal("")
		require.Error(t, err)
		assert.Empty(t, p)
	})

	t.Run("empty id", func(t *testing.T) {
		p, err := jobpayload.Unmarshal(types.MessageIDNil.String())
		require.Error(t, err)
		assert.Empty(t, p)
	})
}
