package internalerrors

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	defaultCode    = 500
	defaultMessage = "something went wrong"
)

// ServerError is used to return custom error codes to client.
type ServerError struct {
	Code    int
	Message string
	cause   error
}

func NewServerError(code int, msg string, err error) *ServerError {
	return &ServerError{Code: code, Message: msg, cause: err}
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %v", s.Message, s.cause)
}

func (s *ServerError) Unwrap() error {
	return s.cause
}

func GetServerErrorCode(err error) int {
	code, _, _ := ProcessServerError(err)
	return code
}

// ProcessServerError tries to retrieve from given error it's code, message and some details.
// For example, that fields can be used to build error response for client.
func ProcessServerError(err error) (code int, msg string, details string) {
	if errServer := new(ServerError); errors.As(err, &errServer) {
		return errServer.Code, errServer.Message, errServer.Error()
	}
	if errHTTP := new(echo.HTTPError); errors.As(err, &errHTTP) {
		return errHTTP.Code, fmt.Sprintf("%v", errHTTP.Message), errHTTP.Error()
	}

	return defaultCode, defaultMessage, err.Error()
}
