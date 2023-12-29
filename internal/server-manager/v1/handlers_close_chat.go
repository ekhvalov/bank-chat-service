package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	closechat "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/close-chat"
)

func (h Handlers) PostCloseChat(eCtx echo.Context, params PostCloseChatParams) error {
	var httpRequest CloseChatRequest
	if err := eCtx.Bind(&httpRequest); err != nil {
		return fmt.Errorf("%w: %v", echo.ErrBadRequest, err)
	}

	req := closechat.Request{ChatID: httpRequest.ChatId, ManagerID: middlewares.MustUserID(eCtx), RequestID: params.XRequestID}

	err := h.closeChatUsecase.Handle(eCtx.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, closechat.ErrNoActiveProblemInChat):
			return internalerrors.NewServerError(int(ErrorCodeNoActiveProblemInChat), err.Error(), err)
		case errors.Is(err, closechat.ErrInvalidRequest):
			return fmt.Errorf("%w: %v", echo.ErrBadRequest, err)
		}
		return fmt.Errorf("%w: %v", echo.ErrInternalServerError, err)
	}

	err = eCtx.JSON(http.StatusOK, CloseChatResponse{
		Data:  new(interface{}),
		Error: nil,
	})
	if err != nil {
		h.lg.Error("closeChat response", zap.Error(err))
	}

	return nil
}
