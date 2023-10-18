package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func NewRecover(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.With(
				zap.Error(err),
				zap.ByteString("stack", stack),
			).Error("panic recovered")
			return err
		},
	})
}
