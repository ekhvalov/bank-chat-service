package clientv1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	gethistory "github.com/ekhvalov/bank-chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/client/send-message"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=clientv1mocks

var _ ServerInterface = (*Handlers)(nil)

type getHistoryUseCase interface {
	Handle(ctx context.Context, req gethistory.Request) (gethistory.Response, error)
}

type sendMessageUseCase interface {
	Handle(ctx context.Context, req sendmessage.Request) (sendmessage.Response, error)
}

//go:generate options-gen -out-filename=handlers_options.gen.go -from-struct=Options
type Options struct {
	logger        *zap.Logger        `option:"mandatory" validate:"required"`
	getHistoryUC  getHistoryUseCase  `option:"mandatory" validate:"required"`
	sendMessageUC sendMessageUseCase `option:"mandatory" validate:"required"`
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := opts.Validate(); err != nil {
		return Handlers{}, fmt.Errorf("validate options: %v", err)
	}
	return Handlers{Options: opts}, nil
}
