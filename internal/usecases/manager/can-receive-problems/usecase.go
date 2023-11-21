package canreceiveproblems

import (
	"context"
	"errors"
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -typed -package=canreceiveproblemsmocks

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrManagerPool        = errors.New("manager pool error")
	ErrManagerLoadService = errors.New("manager load service error")
)

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Contains(ctx context.Context, managerID types.UserID) (bool, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	loadService managerLoadService `option:"mandatory" validate:"required"`
	pool        managerPool        `option:"mandatory" validate:"required"`
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate options: %v", err)
	}
	return UseCase{Options: opts}, nil
}

type UseCase struct {
	Options
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	isMgrInPool, err := u.pool.Contains(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrManagerPool, err)
	}
	if isMgrInPool {
		return Response{Result: false}, nil
	}

	canTakeProblem, err := u.loadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrManagerLoadService, err)
	}
	return Response{Result: canTakeProblem}, nil
}
