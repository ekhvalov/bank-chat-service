package managerv1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	closechat "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/close-chat"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/get-chats"
	sendmessage "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/send-message"
)

var _ ServerInterface = (*Handlers)(nil)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -typed -package=managerv1mocks

type canReceiveProblemsUsecase interface {
	Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error)
}

type freeHandsUsecase interface {
	Handle(ctx context.Context, req freehands.Request) error
}

type getChatsUsecase interface {
	Handle(ctx context.Context, req getchats.Request) (getchats.Response, error)
}

type getChatHistoryUsecase interface {
	Handle(ctx context.Context, req getchathistory.Request) (getchathistory.Response, error)
}

type closeChatUseCase interface {
	Handle(ctx context.Context, req closechat.Request) error
}

type sendMessageUseCase interface {
	Handle(ctx context.Context, req sendmessage.Request) (sendmessage.Response, error)
}

//go:generate options-gen -out-filename=handlers_options.gen.go -from-struct=Options
type Options struct {
	lg                   *zap.Logger               `option:"mandatory" validate:"required"`
	canReceiveProblemsUC canReceiveProblemsUsecase `option:"mandatory" validate:"required"`
	freeHandsUC          freeHandsUsecase          `option:"mandatory" validate:"required"`
	getChatsUC           getChatsUsecase           `option:"mandatory" validate:"required"`
	getChatHistoryUC     getChatHistoryUsecase     `option:"mandatory" validate:"required"`
	sendMessageUsecase   sendMessageUseCase        `option:"mandatory" validate:"required"`
	closeChatUsecase     closeChatUseCase          `option:"mandatory" validate:"required"`
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
