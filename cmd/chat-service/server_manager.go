package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/config"
	problemsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/problems"
	servermanager "github.com/ekhvalov/bank-chat-service/internal/server-manager"
	errhandlermgr "github.com/ekhvalov/bank-chat-service/internal/server-manager/errhandler"
	managerv1 "github.com/ekhvalov/bank-chat-service/internal/server-manager/v1"
	"github.com/ekhvalov/bank-chat-service/internal/server/errhandler"
	managerload "github.com/ekhvalov/bank-chat-service/internal/services/manager-load"
	inmemmanagerpool "github.com/ekhvalov/bank-chat-service/internal/services/manager-pool/in-mem"
	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	canreceiveproblems "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/ekhvalov/bank-chat-service/internal/usecases/manager/free-hands"
)

const (
	logNameServerManager = "server-manager"
)

func initServerManager(
	cfg config.Config,
	storeDB *store.Database,
) (*servermanager.Server, error) {
	lg := zap.L().Named(logNameServerManager)

	v1Swagger, err := managerv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get swagger: %v", err)
	}

	keycloakClient, err := initKeycloakClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create keycloak client: %v", err)
	}

	v1Handlers, err := initV1ManagerHandlers(lg, storeDB, cfg.Services.ManagerLoad)
	if err != nil {
		return nil, fmt.Errorf("create v1Handlers: %v", err)
	}

	errHandler, err := errhandler.New(errhandler.NewOptions(lg, cfg.Global.IsProduction(), errhandlermgr.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create error handler: %v", err)
	}

	server, err := servermanager.New(servermanager.NewOptions(
		cfg.Servers.Manager.Addr,
		cfg.Servers.Manager.AllowOrigins,
		cfg.Servers.Manager.RequiredAccess.Resource,
		cfg.Servers.Manager.RequiredAccess.Role,
		lg,
		v1Swagger,
		v1Handlers,
		keycloakClient,
		errHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return server, nil
}

func initV1ManagerHandlers(lg *zap.Logger, storeDB *store.Database, cfg config.ManagerLoadService) (managerv1.Handlers, error) {
	repo, err := problemsrepo.New(problemsrepo.NewOptions(storeDB))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create problems repo: %v", err)
	}

	pool := inmemmanagerpool.New()

	loadService, err := managerload.New(managerload.NewOptions(cfg.MaxProblemsAtSameTime, repo))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create manager load service: %v", err)
	}

	canReceiveProblemsUsecase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(loadService, pool))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create canReceiveProblems usecase: %v", err)
	}

	freeHandsUsecase, err := freehands.New(freehands.NewOptions(loadService, pool))
	if err != nil {
		return managerv1.Handlers{}, fmt.Errorf("create freeHands usecase: %v", err)
	}

	return managerv1.NewHandlers(managerv1.NewOptions(lg, canReceiveProblemsUsecase, freeHandsUsecase))
}
