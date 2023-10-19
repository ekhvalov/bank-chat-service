package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/errors"
)

const (
	messageSuccess = "success"
	messageError   = "error"
)

func NewLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().Method == echo.OPTIONS
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			level := zap.InfoLevel
			message := messageSuccess
			if v.Status >= http.StatusBadRequest {
				level = zap.ErrorLevel
				message = messageError
			}
			var uid string
			if id, ok := userID(c); ok {
				uid = id.String()
			}
			fields := append(fieldsFromValues(v), zap.String("user_id", uid))
			logger.Log(level, message, fields...)
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
		HandleError:  true,
	})
}

func fieldsFromValues(v middleware.RequestLoggerValues) []zap.Field {
	status := v.Status
	if v.Error != nil {
		status = errors.GetServerErrorCode(v.Error)
	}
	return []zap.Field{
		zap.Duration("latency", v.Latency),
		zap.String("host", v.Host),
		zap.String("method", v.Method),
		zap.String("path", v.URIPath),
		zap.String("request_id", v.RequestID),
		zap.String("user_agent", v.UserAgent),
		zap.Int("status", status),
		zap.Error(v.Error),
	}
}
