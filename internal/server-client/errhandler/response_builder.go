package errhandler

import (
	clientv1 "github.com/ekhvalov/bank-chat-service/internal/server-client/v1"
	"github.com/ekhvalov/bank-chat-service/pkg/pointer"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	return Response{Error: clientv1.Error{
		Code:    clientv1.ErrorCode(code),
		Details: pointer.PtrWithZeroAsNil(details),
		Message: msg,
	}}
}