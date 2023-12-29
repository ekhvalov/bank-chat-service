package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/send-message"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	var httpRequest SendMessageRequest
	if err := eCtx.Bind(&httpRequest); err != nil {
		return fmt.Errorf("%w: %v", echo.ErrBadRequest, err)
	}

	request := sendmessage.Request{
		ID:          params.XRequestID,
		ManagerID:   managerID,
		ChatID:      httpRequest.ChatId,
		MessageBody: httpRequest.MessageBody,
	}
	response, err := h.sendMessageUsecase.Handle(ctx, request)
	if err != nil {
		if errors.Is(err, sendmessage.ErrInvalidRequest) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = eCtx.JSON(http.StatusOK, SendMessageResponse{
		Data: &MessageWithoutBody{
			AuthorId:  managerID,
			CreatedAt: response.CreatedAt,
			Id:        response.MessageID,
		},
		Error: nil,
	})
	if err != nil {
		h.lg.Error("SendMessageResponse", zap.Error(err))
	}

	return nil
}
