package freehands

import (
	"context"
	"errors"
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -typed -package=freehandsmocks

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrManagerOverloaded = errors.New("manager cannot take more problems")
)

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Contains(ctx context.Context, managerID types.UserID) (bool, error)
	Put(ctx context.Context, managerID types.UserID) error
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

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	isMgrInPool, err := u.pool.Contains(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("pool contains: %v", err)
	}
	if isMgrInPool {
		return nil
	}

	canTakeProblem, err := u.loadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("loadService can take: %v", err)
	}
	if !canTakeProblem {
		return ErrManagerOverloaded
	}

	if err = u.pool.Put(ctx, req.ManagerID); err != nil {
		return fmt.Errorf("put manager to pool: %v", err)
	}
	return nil
}
