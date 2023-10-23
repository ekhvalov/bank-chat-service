package errhandler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalerrors "github.com/ekhvalov/bank-chat-service/internal/errors"
)

var _ echo.HTTPErrorHandler = Handler{}.Handle

//go:generate options-gen -out-filename=errhandler_options.gen.go -from-struct=Options
type Options struct {
	logger          *zap.Logger                                    `option:"mandatory" validate:"required"`
	productionMode  bool                                           `option:"mandatory"`
	responseBuilder func(code int, msg string, details string) any `option:"mandatory" validate:"required"`
}

type Handler struct {
	lg              *zap.Logger
	productionMode  bool
	responseBuilder func(code int, msg string, details string) any
}

func New(opts Options) (Handler, error) {
	if err := opts.Validate(); err != nil {
		return Handler{}, fmt.Errorf("options validate: %v", err)
	}
	return Handler{lg: opts.logger, productionMode: opts.productionMode, responseBuilder: opts.responseBuilder}, nil
}

func (h Handler) Handle(err error, eCtx echo.Context) {
	code, msg, details := internalerrors.ProcessServerError(err)
	if h.productionMode {
		details = ""
	}
	if errSend := eCtx.JSON(http.StatusOK, h.responseBuilder(code, msg, details)); errSend != nil {
		h.lg.Error("respond error", zap.Error(errSend))
	}
}
