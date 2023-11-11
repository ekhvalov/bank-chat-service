package errhandler

import (
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

type Response struct {
	Error managerv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	return Response{Error: managerv1.Error{
		Code:    managerv1.ErrorCode(code),
		Details: pointer.PtrWithZeroAsNil(details),
		Message: msg,
	}}
}
