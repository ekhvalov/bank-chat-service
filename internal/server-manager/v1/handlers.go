package managerv1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
)

var _ ServerInterface = (*Handlers)(nil)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=managerv1mocks

type canReceiveProblemsUsecase interface {
	Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error)
}

type freeHandsUsecase interface {
	Handle(ctx context.Context, req freehands.Request) error
}

//go:generate options-gen -out-filename=handlers_options.gen.go -from-struct=Options
type Options struct {
	lg                   *zap.Logger               `option:"mandatory" validate:"required"`
	canReceiveProblemsUC canReceiveProblemsUsecase `option:"mandatory" validate:"required"`
	freeHandsUC          freeHandsUsecase          `option:"mandatory" validate:"required"`
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
