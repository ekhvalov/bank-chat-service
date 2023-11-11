package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
	"github.com/ekhvalov/bank-chat-service/internal/middlewares"
	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(eCtx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
	req := canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: middlewares.MustUserID(eCtx),
	}

	result, err := h.canReceiveProblemsUC.Handle(eCtx.Request().Context(), req)
	if err != nil {
		if errors.Is(err, canreceiveproblems.ErrInvalidRequest) {
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		return fmt.Errorf("%w: %v", echo.ErrInternalServerError, err)
	}

	err = eCtx.JSON(http.StatusOK, GetFreeHandsBtnAvailabilityResponse{
		Data:  &FreeHandsAvailability{Available: result.Result},
		Error: nil,
	})
	if err != nil {
		h.lg.Error("postGetFreeHandsBtnAvailability response", zap.Error(err))
	}
	return nil
}

func (h Handlers) PostFreeHands(eCtx echo.Context, params PostFreeHandsParams) error {
	req := freehands.Request{
		ID:        params.XRequestID,
		ManagerID: middlewares.MustUserID(eCtx),
	}
	err := h.freeHandsUC.Handle(eCtx.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, freehands.ErrInvalidRequest):
			return internalerrors.NewServerError(http.StatusBadRequest, "invalid request", err)
		case errors.Is(err, freehands.ErrManagerOverloaded):
			return internalerrors.NewServerError(int(ErrorCodeManagerOverloadedError), "manager overloaded", err)
		}
		return fmt.Errorf("%w: %v", echo.ErrInternalServerError, err)
	}

	err = eCtx.JSON(http.StatusOK, FreeHandsResponse{Data: new(interface{})})
	if err != nil {
		h.lg.Error("postFreeHands response", zap.Error(err))
	}
	return nil
}
