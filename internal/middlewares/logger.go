package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func NewLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().Method == echo.OPTIONS
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var uid string
			if id, ok := userID(c); ok {
				uid = id.String()
			}
			logger.Info(
				"success",
				zap.Duration("latency", v.Latency),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("path", v.URIPath),
				zap.String("request_id", v.RequestID),
				zap.String("user_agent", v.UserAgent),
				zap.Int("status", v.Status),
				zap.String("user_id", uid),
				zap.Error(v.Error),
			)
			return nil
		},
		LogLatency:   true,
		LogHost:      true,
		LogMethod:    true,
		LogURIPath:   true,
		LogRequestID: true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
	})
}
