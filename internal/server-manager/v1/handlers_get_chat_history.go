package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	getchathistory "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chat-history"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

func (h Handlers) PostGetChatHistory(eCtx echo.Context, params PostGetChatHistoryParams) error {
	var req GetChatHistoryRequest
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("%w: %v", echo.ErrBadRequest, err)
	}
	request := getchathistory.Request{
		ID:        params.XRequestID,
		ManagerID: middlewares.MustUserID(eCtx),
		ChatID:    req.ChatId,
		PageSize:  pointer.Indirect(req.PageSize),
		Cursor:    pointer.Indirect(req.Cursor),
	}

	result, err := h.getChatHistoryUC.Handle(eCtx.Request().Context(), request)
	if err != nil {
		if errors.Is(err, getchathistory.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		return fmt.Errorf("%w: %v", echo.ErrInternalServerError, err)
	}

	messages := make([]Message, len(result.Messages))
	for i, msg := range result.Messages {
		messages[i] = Message{
			AuthorId:  msg.AuthorID,
			Body:      msg.Body,
			CreatedAt: msg.CreatedAt,
			Id:        msg.ID,
		}
	}
	err = eCtx.JSON(http.StatusOK, GetChatHistoryResponse{
		Data:  &MessagesPage{Messages: messages},
		Error: nil,
	})
	if err != nil {
		h.lg.Error("postGetChatHistory response", zap.Error(err))
	}
	return nil
}
