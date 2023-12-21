package managerv1_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	closechat "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/close-chat"
)

func (s *HandlersSuite) TestPostCloseChat_BadRequestError() {
	// Arrange.
	reqID := types.NewRequestID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", `{"incomplete":"request`)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	var errHTTP *echo.HTTPError
	s.Require().ErrorAs(err, &errHTTP)
	s.Equal(http.StatusBadRequest, errHTTP.Code)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestPostCloseChat_Usecase_InvalidRequestError() {
	// Arrange.
	reqID := types.NewRequestID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", "{}")

	s.closeChatUsecase.EXPECT().Handle(eCtx.Request().Context(), gomock.Any()).Return(closechat.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	var errHTTP *echo.HTTPError
	s.Require().ErrorAs(err, &errHTTP)
	s.Equal(http.StatusBadRequest, errHTTP.Code)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestPostCloseChat_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	httpReq := fmt.Sprintf(`{"chatId": %q}`, chatID)
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", httpReq)

	s.closeChatUsecase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ChatID:    chatID,
		ManagerID: s.managerID,
		RequestID: reqID,
	}).Return(errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	var errHTTP *echo.HTTPError
	s.Require().ErrorAs(err, &errHTTP)
	s.Equal(http.StatusInternalServerError, errHTTP.Code)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestPostCloseChat_Usecase_NoActiveProblemInChatError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	httpReq := fmt.Sprintf(`{"chatId": %q}`, chatID)
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", httpReq)

	s.closeChatUsecase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ChatID:    chatID,
		ManagerID: s.managerID,
		RequestID: reqID,
	}).Return(closechat.ErrNoActiveProblemInChat)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	var errServer *internalerrors.ServerError
	s.Require().ErrorAs(err, &errServer)
	s.EqualValues(managerv1.ErrorCodeNoActiveProblemInChat, errServer.Code)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestPostCloseChat_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()

	httpReq := fmt.Sprintf(`{"chatId": %q}`, chatID)
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", httpReq)

	s.closeChatUsecase.EXPECT().Handle(eCtx.Request().Context(), closechat.Request{
		ChatID:    chatID,
		ManagerID: s.managerID,
		RequestID: reqID,
	}).Return(nil)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
	"data": null
}`, resp.Body.String())
}
