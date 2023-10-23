package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	gethistory "github.com/ekhvalov/bank-chat-service/internal/usecases/client/get-history"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)
	var req GetHistoryRequest
	if err := eCtx.Bind(&req); err != nil {
		return fmt.Errorf("%w: %v", echo.ErrBadRequest, err)
	}

	request := gethistory.Request{
		ID:       params.XRequestID,
		ClientID: clientID,
		PageSize: pointer.Indirect(req.PageSize),
		Cursor:   pointer.Indirect(req.Cursor),
	}
	response, err := h.getHistoryUC.Handle(ctx, request)
	if err != nil {
		switch {
		case errors.Is(err, gethistory.ErrInvalidRequest):
			fallthrough
		case errors.Is(err, gethistory.ErrInvalidCursor):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return fmt.Errorf("%w: %w", echo.ErrInternalServerError, err)
	}

	messages := make([]Message, 0, len(response.Messages))
	for _, m := range response.Messages {
		messages = append(messages, Message{
			AuthorId:   pointer.PtrWithZeroAsNil(m.AuthorID),
			Body:       m.Body,
			CreatedAt:  m.CreatedAt,
			Id:         m.ID,
			IsBlocked:  m.IsBlocked,
			IsReceived: m.IsReceived,
			IsService:  m.IsService,
		})
	}
	err = eCtx.JSON(http.StatusOK, GetHistoryResponse{
		Data:  &MessagesPage{Messages: messages, Next: response.NextCursor},
		Error: nil,
	})
	if err != nil {
		h.logger.Error("getHistory response", zap.Error(err))
	}
	return nil
}
