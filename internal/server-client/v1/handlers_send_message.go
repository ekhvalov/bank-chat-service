package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	errs "github.com/ekhvalov/bank-chat-service/internal/errors"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/client/send-message"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)
	var httpRequest SendMessageRequest
	if err := eCtx.Bind(&httpRequest); err != nil {
		return fmt.Errorf("%w: %v", echo.ErrBadRequest, err)
	}

	usecaseRequest := sendmessage.Request{
		ID:          params.XRequestID,
		ClientID:    clientID,
		MessageBody: httpRequest.MessageBody,
	}
	response, err := h.sendMessageUC.Handle(ctx, usecaseRequest)
	if err != nil {
		switch {
		case errors.Is(err, sendmessage.ErrInvalidRequest):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case errors.Is(err, sendmessage.ErrChatNotCreated):
			return errs.NewServerError(int(ErrorCodeCreateChatError), "create chat error", err)
		case errors.Is(err, sendmessage.ErrProblemNotCreated):
			return errs.NewServerError(int(ErrorCodeCreateProblemError), "create problem error", err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = eCtx.JSON(http.StatusOK, SendMessageResponse{
		Data: &MessageHeader{
			AuthorID:  pointer.PtrWithZeroAsNil(clientID),
			CreatedAt: response.CreatedAt,
			ID:        response.MessageID,
		},
		Error: nil,
	})
	if err != nil {
		h.logger.Error("SendMessageResponse", zap.Error(err))
	}

	return nil
}
