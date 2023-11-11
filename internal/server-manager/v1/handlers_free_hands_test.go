package managerv1_test

import (
	"errors"
	"net/http"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
	"github.com/ekhvalov/bank-chat-service/internal/types"
	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
)

func (s *HandlersSuite) TestGetFreeHandsBtnAvailability_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.canReceiveProblemsUsecase.EXPECT().Handle(eCtx.Request().Context(), canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(canreceiveproblems.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetFreeHandsBtnAvailability(eCtx, managerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetFreeHandsBtnAvailability_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.canReceiveProblemsUsecase.EXPECT().Handle(eCtx.Request().Context(), canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(canreceiveproblems.Response{Result: true}, nil)

	// Action.
	err := s.handlers.PostGetFreeHandsBtnAvailability(eCtx, managerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data":
    {
        "available": true
    }
}`, resp.Body.String())
}

func (s *HandlersSuite) Test_PostFreeHands_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	s.freeHandsUsecase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	})

	// Action
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
	"data": null
}`, resp.Body.String())
}

func (s *HandlersSuite) Test_PostFreeHands_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	errUsecase := errors.New("usecase error")
	s.freeHandsUsecase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(errUsecase)

	// Action
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert
	s.Require().Error(err)
	s.Equal(http.StatusInternalServerError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) Test_PostFreeHands_Usecase_InvalidRequest_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	s.freeHandsUsecase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(freehands.ErrInvalidRequest)

	// Action
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert
	s.Require().ErrorIs(err, freehands.ErrInvalidRequest)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) Test_PostFreeHands_Usecase_ManagerOverloaded_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	s.freeHandsUsecase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(freehands.ErrManagerOverloaded)

	// Action
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert
	s.Require().ErrorIs(err, freehands.ErrManagerOverloaded)
	s.EqualValues(managerv1.ErrorCodeManagerOverloadedError, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}
