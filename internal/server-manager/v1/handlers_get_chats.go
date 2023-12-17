package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	getchats "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chats"
)

func (h Handlers) PostGetChats(eCtx echo.Context, params PostGetChatsParams) error {
	req := getchats.Request{
		ID:        params.XRequestID,
		ManagerID: middlewares.MustUserID(eCtx),
	}

	result, err := h.getChatsUC.Handle(eCtx.Request().Context(), req)
	if err != nil {
		if errors.Is(err, getchats.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		return fmt.Errorf("%w: %v", echo.ErrInternalServerError, err)
	}

	chats := make([]Chat, len(result.Chats))
	for i, chat := range result.Chats {
		chats[i] = Chat{ChatId: chat.ID, ClientId: chat.ClientID}
	}
	err = eCtx.JSON(http.StatusOK, GetChatsResponse{
		Data:  &ChatList{Chats: chats},
		Error: nil,
	})
	if err != nil {
		h.lg.Error("postGetChats response", zap.Error(err))
	}
	return nil
}
