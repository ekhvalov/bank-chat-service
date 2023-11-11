package errhandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ekhvalov/bank-chat-service/internal/server-manager/errhandler"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
)

func TestResponseBuilder(t *testing.T) {
	t.Run("with details", func(t *testing.T) {
		err := errhandler.ResponseBuilder(5000, "hello", "world")

		resp, ok := err.(errhandler.Response)
		require.True(t, ok)
		require.IsType(t, managerv1.Error{}, resp.Error)

		assert.Equal(t, managerv1.ErrorCodeManagerOverloadedError, resp.Error.Code)
		assert.Equal(t, "hello", resp.Error.Message)
		require.NotNil(t, resp.Error.Details)
		assert.Equal(t, "world", *resp.Error.Details)
	})

	t.Run("without details", func(t *testing.T) {
		err := errhandler.ResponseBuilder(5000, "hello", "")

		resp, ok := err.(errhandler.Response)
		require.True(t, ok)
		require.IsType(t, managerv1.Error{}, resp.Error)

		assert.Equal(t, managerv1.ErrorCodeManagerOverloadedError, resp.Error.Code)
		assert.Equal(t, "hello", resp.Error.Message)
		assert.Nil(t, resp.Error.Details)
	})
}
